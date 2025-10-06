package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
	"github.com/cprakhar/relief-ops/shared/tools"
)

type KafkaConfig struct {
	Brokers string
	GroupID string
}

// KafkaClient wraps the Kafka producer and consumer
type KafkaClient struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

// NewKafkaClient initializes a new Kafka client
func NewKafkaClient(serviceName string, cfg *KafkaConfig) (*KafkaClient, error) {
	producerCfg := &kafka.ConfigMap{
		"bootstrap.servers":  cfg.Brokers,
		"acks":               "all",
		"compression.type":   "snappy",
		"client.id":          serviceName,
		"security.protocol":  "PLAINTEXT",
		"enable.idempotence": true,
	}
	producer, err := kafka.NewProducer(producerCfg)
	if err != nil {
		return nil, err
	}

	consumerCfg := &kafka.ConfigMap{
		"bootstrap.servers":        cfg.Brokers,
		"group.id":                 cfg.GroupID,
		"auto.offset.reset":        "earliest",
		"enable.auto.commit":       false,
		"session.timeout.ms":       60000,
		"allow.auto.create.topics": true,
		"security.protocol":        "PLAINTEXT",
		"client.id":                serviceName,
	}
	consumer, err := kafka.NewConsumer(consumerCfg)
	if err != nil {
		producer.Close()
		return nil, err
	}

	return &KafkaClient{
		Producer: producer,
		Consumer: consumer,
	}, nil
}

// Produce sends a message to the specified topic with retries and error handling
func (kc *KafkaClient) Produce(ctx context.Context, topic string, key string, value []byte) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}

	// Use delivery channel to wait for confirmation
	deliveryChan := make(chan kafka.Event, 1)

	// Send the message with tracing
	err := traces.TracedProduce(ctx, kc.Producer.Produce, message, deliveryChan)
	if err != nil {
		return err
	}

	// Wait for delivery confirmation or context cancellation
	select {
	case event := <-deliveryChan:
		switch ev := event.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				return ev.TopicPartition.Error
			}
			// Message delivered successfully
			return nil
		default:
			return kafka.NewError(kafka.ErrUnknown, "unexpected delivery event", false)
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}

// MessageHandler defines the function signature for handling messages
type MessageHandler func(ctx context.Context, event, key string, value []byte) error

// Consume starts consuming messages from the specified topics and processes them using the provided handler
func (kc *KafkaClient) Consume(ctx context.Context, topics []string, handler MessageHandler) error {
	// Subscribe to topics
	if err := kc.Consumer.SubscribeTopics(topics, nil); err != nil {
		return err
	}

	// Ensure consumer is closed on context cancellation
	go func() {
		<-ctx.Done()
		kc.Consumer.Close()
	}()

	// Poll for messages
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		event := kc.Consumer.Poll(100)
		if event == nil {
			continue
		}

		switch ev := event.(type) {
		case *kafka.Message:
			// Retry handling the message with exponential backoff
			// If it fails after retries, send to DLQ
			retryConfig := &tools.RetryConfig{
				MaxAttempts:   5,
				InitialDelay:  100 * time.Millisecond,
				MaxDelay:      5 * time.Second,
				BackoffFactor: 2.0,
				Jitter:        true,
			}

			if err := tools.RetryWithBackoff(ctx, retryConfig, func() error {
				// Call the traced handler
				return traces.TracedConsumeHandler(ctx, handler, ev)
			}); err != nil {
				log.Printf("Error handling message after retries: %v", err)
				// Send to Dead Letter Queue (DLQ)
				if dlqErr := kc.DLQ(ctx, *ev.TopicPartition.Topic, string(ev.Key), ev.Value); dlqErr != nil {
					log.Printf("CRITICAL: Failed to send message to DLQ: %v", dlqErr)
					continue
				}
				log.Printf("Message sent to DLQ: topic=%s, key=%s", fmt.Sprintf("%s-DLQ", *ev.TopicPartition.Topic), string(ev.Key))
			}
			// Commit the message offset after processing or sending to DLQ
			_, err := kc.Consumer.CommitMessage(ev)
			if err != nil {
				log.Printf("Error committing message: %v", err)
			}
		case kafka.Error:
			log.Printf("Kafka error: %v", ev)
		default:
			// Ignore other event types
			log.Printf("Ignored event: %v", ev)
		}
	}
}

// DLQ sends the message to a Dead Letter Queue topic
func (kc *KafkaClient) DLQ(ctx context.Context, topic string, key string, value []byte) error {
	dlqTopic := topic + "-DLQ"
	dlqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return kc.Produce(dlqCtx, dlqTopic, key, value)
}

// Close cleans up the Kafka producer and consumer
func (kc *KafkaClient) Close() {
	kc.Producer.Close()
	kc.Consumer.Close()
}
