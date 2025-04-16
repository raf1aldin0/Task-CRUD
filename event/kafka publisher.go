package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Publisher interface {
	Publish(ctx context.Context, key, value string) error
	Close() error
}

type kafkaPublisher struct {
	writer *kafka.Writer
	tracer trace.Tracer
}

// NewPublisher menginisialisasi publisher Kafka
func NewPublisher(cfg KafkaConfig) Publisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{}, // Load balancing antar partition
		RequiredAcks: kafka.RequireAll,    // Tunggu semua broker ack
		Async:        false,               // Sync publish, lebih aman untuk transaksi
	}

	return &kafkaPublisher{
		writer: writer,
		tracer: otel.Tracer("kafka-publisher"),
	}
}

// Publish mengirimkan pesan ke Kafka dengan tracing
func (kp *kafkaPublisher) Publish(ctx context.Context, key, value string) error {
	ctx, span := kp.tracer.Start(ctx, "PublishKafkaMessage", trace.WithAttributes(
		attribute.String("message.key", key),
		attribute.String("message.value", value),
		attribute.String("topic", kp.writer.Topic),
	))
	defer span.End()

	message := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
	}

	if err := kp.writer.WriteMessages(ctx, message); err != nil {
		log.Printf("❌ Failed to publish message: %v", err)
		return err
	}

	log.Printf("📤 Published message: Key=%s, Value=%s", key, value)
	return nil
}

// Close menutup Kafka writer
func (kp *kafkaPublisher) Close() error {
	return kp.writer.Close()
}
