package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Consumer interface {
	Start(ctx context.Context) error
	Close() error
}

type kafkaConsumer struct {
	reader *kafka.Reader
	tracer trace.Tracer
}

// ✅ Auto Commit aktif: CommitInterval > 0
func NewConsumer(cfg KafkaConfig) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        "task-crud-group",
		MinBytes:       10e3,             // 10KB
		MaxBytes:       10e6,             // 10MB
		CommitInterval: 5 * time.Second,  // ✅ Auto commit setiap 5 detik
	})

	return &kafkaConsumer{
		reader: reader,
		tracer: otel.Tracer("kafka-consumer"),
	}
}

func (kc *kafkaConsumer) Start(ctx context.Context) error {
	fmt.Println("🟢 Kafka consumer listening (auto commit)...")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		m, err := kc.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				log.Println("⚠️ Kafka consumer stopped gracefully")
				return nil
			}
			log.Printf("❌ Kafka consumer error: %v, retrying in 2 seconds...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Tracing per message
		_, span := kc.tracer.Start(ctx, "ConsumeKafkaMessage", trace.WithAttributes(
			attribute.String("message.key", string(m.Key)),
			attribute.String("message.value", string(m.Value)),
		))
		defer span.End()

		log.Printf("📥 Kafka received message: Key=%s, Value=%s\n", string(m.Key), string(m.Value))

		// ⚠️ Tidak perlu manual commit, Kafka akan commit otomatis
	}
}

func (kc *kafkaConsumer) Close() error {
	return kc.reader.Close()
}
