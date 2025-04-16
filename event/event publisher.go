package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	
	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type EventPublisher interface {
	PublishRepositoryCreated(ctx context.Context, data interface{}) error
	PublishUserUpdated(ctx context.Context, data interface{}) error
	Close() error
}

type eventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher(cfg KafkaConfig) EventPublisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic, // default topic, bisa di override per-event jika mau
		Balancer: &kafka.LeastBytes{},
	}

	return &eventPublisher{
		writer: writer,
	}
}

func (ep *eventPublisher) publish(ctx context.Context, key string, eventType string, data interface{}) error {
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(eventType)},
		},
	}

	if err := ep.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("❌ Failed to publish %s: %v", eventType, err)
		return err
	}

	log.Printf("📤 Published event [%s]: %s", eventType, string(value))
	return nil
}

func (ep *eventPublisher) PublishRepositoryCreated(ctx context.Context, data interface{}) error {
	return ep.publish(ctx, "repository", "RepositoryCreated", data)
}

func (ep *eventPublisher) PublishUserUpdated(ctx context.Context, data interface{}) error {
	return ep.publish(ctx, "user", "UserUpdated", data)
}

func (ep *eventPublisher) Close() error {
	return ep.writer.Close()
}
