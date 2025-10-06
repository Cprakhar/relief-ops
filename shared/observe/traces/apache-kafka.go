package traces

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type kafkaHeaderCarrier []kafka.Header

func (c *kafkaHeaderCarrier) Get(key string) string {
	for _, h := range *c {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *kafkaHeaderCarrier) Set(key, value string) {
	for i, h := range *c {
		if h.Key == key {
			(*c)[i].Value = []byte(value)
			return
		}
	}
	*c = append(*c, kafka.Header{Key: key, Value: []byte(value)})
}

func (c *kafkaHeaderCarrier) Keys() []string {
	keys := make([]string, len(*c))
	for i, h := range *c {
		keys[i] = h.Key
	}
	return keys
}

// tracedProduce wraps the Produce method with OpenTelemetry tracing.
func TracedProduce(ctx context.Context, produceFunc func(msg *kafka.Message, deliveryChan chan kafka.Event) error, msg *kafka.Message, deliveryChan chan kafka.Event) error {
	tracer := otel.GetTracerProvider().Tracer("apache-kafka")

	ctx, span := tracer.Start(ctx, "apache-kafka.produce",
		trace.WithAttributes(
			attribute.String("messaging.system", "apache-kafka"),
			attribute.String("messaging.destination", *msg.TopicPartition.Topic),
			attribute.String("messaging.operation", "produce"),
		),
	)
	defer span.End()

	span.SetAttributes(attribute.String("messaging.apache-kafka.message_key", string(msg.Key)))

	if msg.Headers == nil {
		msg.Headers = []kafka.Header{}
	}

	carrier := kafkaHeaderCarrier(msg.Headers)
	otel.GetTextMapPropagator().Inject(ctx, &carrier)
	msg.Headers = []kafka.Header(carrier)

	err := produceFunc(msg, deliveryChan)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "Message produced successfully")
	return nil
}

func TracedConsumeHandler(ctx context.Context, handler func(ctx context.Context, topic, key string, value []byte) error, msg *kafka.Message) error {
	var carrier kafkaHeaderCarrier
	if msg.Headers != nil {
		carrier = kafkaHeaderCarrier(msg.Headers)
	}

	ctx = otel.GetTextMapPropagator().Extract(ctx, &carrier)

	tracer := otel.GetTracerProvider().Tracer("apache-kafka")
	ctx, span := tracer.Start(ctx, "apache-kafka.consume",
		trace.WithAttributes(
			attribute.String("messaging.system", "apache-kafka"),
			attribute.String("messaging.destination", *msg.TopicPartition.Topic),
			attribute.String("messaging.operation", "consume"),
			attribute.String("messaging.apache-kafka.message_key", string(msg.Key)),
			attribute.Int64("messaging.apache-kafka.partition", int64(msg.TopicPartition.Partition)),
			attribute.Int64("messaging.apache-kafka.offset", int64(msg.TopicPartition.Offset)),
		),
	)
	defer span.End()

	err := handler(ctx, *msg.TopicPartition.Topic, string(msg.Key), msg.Value)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "Message consumed and processed successfully")
	return nil
}
