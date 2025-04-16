package kafka

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// LoadKafkaConfig membaca konfigurasi Kafka dari environment variable
func LoadKafkaConfig() KafkaConfig {
	// Bisa menggunakan banyak broker, dipisahkan dengan koma
	brokerEnv := os.Getenv("KAFKA_BROKERS")
	if brokerEnv == "" {
		brokerEnv = "kafka:9092" // Default jika tidak di-set
	}
	brokers := strings.Split(brokerEnv, ",")

	// Baca topic
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "repository-topic"
	}

	return KafkaConfig{
		Brokers: brokers,
		Topic:   topic,
	}
}
