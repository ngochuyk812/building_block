package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ngochuyk812/building_block/infrastructure/eventbus"
	"github.com/ngochuyk812/building_block/infrastructure/helpers"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

func NewProceduer(brokerUrls []string, topic string) (eventbus.Producer, error) {
	instance := &kafkaWriter{}
	err := instance.configure(brokerUrls, topic)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

type kafkaWriter struct {
	writer *kafka.Writer
}

func (k *kafkaWriter) Publish(ctx context.Context, event eventbus.IntegrationEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}
	fmt.Printf("Publish event %s: %s", event.Key(), string(value))
	auth, _ := helpers.AuthContext(ctx)
	authBytes, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("failed to marshal auth context: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(event.Key()),
		Value: value,
		Time:  time.Now(),
		Headers: []kafka.Header{
			{
				Key:   "AuthContext",
				Value: []byte(authBytes),
			},
		},
	}
	return k.writer.WriteMessages(ctx, message)
}

func (k *kafkaWriter) configure(kafkaBrokerUrls []string, topic string) (err error) {
	dialer := &kafka.Dialer{
		Timeout: 10 * time.Second,
	}

	config := kafka.WriterConfig{
		Brokers:          kafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	w := kafka.NewWriter(config)
	k.writer = w

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = k.Publish(ctx, &healthyCheckEvent{
		Value: "ping",
	})
	if err != nil {
		return fmt.Errorf("Kafka connection test failed: %w", err)
	}
	return nil
}

type healthyCheckEvent struct {
	Value string
}

func (h *healthyCheckEvent) Key() string {
	return "health-check"
}
