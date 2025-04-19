package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ngochuyk812/building_block/infrastructure/eventbus"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

func NewProceduer(brokerUrls []string, clientId string, topic string) (eventbus.Producer, error) {
	instance := &kafkaWriter{}
	err := instance.configure(brokerUrls, clientId, topic)
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
	message := kafka.Message{
		Key:   []byte(event.Key()),
		Value: value,
		Time:  time.Now(),
	}
	return k.writer.WriteMessages(ctx, message)
}

func (k *kafkaWriter) configure(kafkaBrokerUrls []string, clientId string, topic string) (err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
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
	return nil
}
