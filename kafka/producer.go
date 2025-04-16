package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Producer interface {
	Publish(ctx context.Context, key string, message interface{}) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
	topic  string
	tracer trace.Tracer
}

// NewProducer mengembalikan instance Kafka producer dengan konfigurasi
func NewProducer(cfg KafkaConfig) Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &kafkaProducer{
		writer: writer,
		topic:  cfg.Topic,
		tracer: otel.Tracer("kafka-producer"),
	}
}

// Publish mengirimkan pesan ke Kafka dan mencatat tracing informasi
func (kp *kafkaProducer) Publish(ctx context.Context, key string, message interface{}) error {
	// Start tracing span untuk operasi publish
	ctx, span := kp.tracer.Start(ctx, "KafkaPublish")
	defer span.End()

	// Encode message ke JSON
	data, err := json.Marshal(message)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Persiapkan Kafka message
	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	}

	// Menambahkan atribut tracing untuk memperkaya informasi
	span.SetAttributes(
		attribute.String("kafka.topic", kp.topic),
		attribute.String("kafka.key", key),
		attribute.Int("message.size", len(data)),
	)

	// Retry logic jika terjadi error saat menulis pesan
	var retryAttempts = 3
	for i := 0; i < retryAttempts; i++ {
		// Publish pesan ke Kafka
		err = kp.writer.WriteMessages(ctx, msg)
		if err != nil {
			log.Printf("❌ Kafka publish error (attempt %d/%d): %v", i+1, retryAttempts, err)
			span.RecordError(err)

			// Jika sudah mencapai jumlah retry maksimal, berhenti
			if i == retryAttempts-1 {
				return fmt.Errorf("failed to write message to Kafka after %d attempts: %w", retryAttempts, err)
			}

			// Tunggu beberapa saat sebelum mencoba lagi
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		// Jika berhasil publish pesan
		log.Println("✅ Kafka message published successfully")
		return nil
	}

	return fmt.Errorf("failed to publish message to Kafka after %d retries", retryAttempts)
}

// Close menutup Kafka writer dengan baik
func (kp *kafkaProducer) Close() error {
	return kp.writer.Close()
}
