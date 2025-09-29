package messaging

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
		"group.protocol":           "CONSUMER",
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

func (kc *KafkaClient) Produce(topic string, key string, value []byte) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}
	return kc.Producer.Produce(message, nil)
}

type MessageHandler func(event, key string, value []byte) error

func (kc *KafkaClient) Consume(ctx context.Context, topics []string, handler MessageHandler) error {
	if err := kc.Consumer.SubscribeTopics(topics, nil); err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		kc.Consumer.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		msg, err := kc.Consumer.ReadMessage(100)
		if err.(kafka.Error).IsTimeout() {
			continue
		}

		if err := handler(*msg.TopicPartition.Topic, string(msg.Key), msg.Value); err != nil {
			log.Printf("Error handling message: %v", err)
		}

		_, err = kc.Consumer.CommitMessage(msg)
		if err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}
}

func (kc *KafkaClient) Close() {
	kc.Producer.Close()
	kc.Consumer.Close()
}
