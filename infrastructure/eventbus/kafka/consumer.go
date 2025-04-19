package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ngochuyk812/building_block/infrastructure/eventbus"
	"github.com/ngochuyk812/building_block/infrastructure/helpers"
	auth_context "github.com/ngochuyk812/building_block/pkg/auth"
	"github.com/segmentio/kafka-go"
)

var _ eventbus.Consumer = (*kafkaConsumer)(nil)

type kafkaConsumer struct {
	reader   *kafka.Reader
	handler  map[string]eventbus.IntegrationEventHandler
	topic    string
	brokers  string
	group_id string
}

func NewConsumer(brokers, topic, group_id string) (eventbus.Consumer, error) {
	instance := &kafkaConsumer{
		handler:  make(map[string]eventbus.IntegrationEventHandler),
		topic:    topic,
		brokers:  brokers,
		group_id: group_id,
	}
	return instance, nil
}

func (k *kafkaConsumer) RegisterHandler(handler eventbus.IntegrationEventHandler) (err error) {
	key := handler.NewEvent().Key()
	log.Printf("Register handler key %s handler %v", key, handler)
	k.handler[key] = handler
	return nil
}

func (k *kafkaConsumer) Run() {
	log.Printf("Connect kafka %s topic %s", k.brokers, k.topic)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	config := kafka.ReaderConfig{
		Brokers:         strings.Split(k.brokers, ","),
		Topic:           k.topic,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
		GroupID:         k.group_id,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		key := string(m.Key)

		if handler, ok := k.handler[key]; ok {
			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))

			event := handler.NewEvent()
			if err := json.Unmarshal(m.Value, event); err != nil {
				log.Printf("unmarshal error for key %s: %v", key, err)
				continue
			}

			ctx := context.Background()

			ctx = AuthContextMiddleware(ctx, m)

			if err := handler.Handle(ctx, event); err != nil {
				log.Printf("handle error for key %s: %v", key, err)
				continue
			}
			if err := reader.CommitMessages(context.Background(), m); err != nil {
				log.Printf("commit error: %v", err)
			}
		}
	}
}

func AuthContextMiddleware(ctx context.Context, msg kafka.Message) context.Context {
	for _, header := range msg.Headers {
		if header.Key == "AuthContext" {
			var auth auth_context.AuthContext
			if err := json.Unmarshal(header.Value, &auth); err == nil {
				return helpers.NewContext(ctx, helpers.AuthContextKey, &auth)
			}
		}
	}
	return ctx
}
