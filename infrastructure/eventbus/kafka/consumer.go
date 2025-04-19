package kafka

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ngochuyk812/building_block/infrastructure/eventbus"
	"github.com/segmentio/kafka-go"
)

var _ eventbus.Consumer = (*kafkaConsumer)(nil)

var (
	kafkaBrokerUrl string
	kafkaTopic     string
)

type kafkaConsumer struct {
	reader  *kafka.Reader
	handler map[string]eventbus.IntegrationEventHandler
}

func (k *kafkaConsumer) RegisterHandler(handler eventbus.IntegrationEventHandler) (err error) {
	key := handler.NewEvent().Key()
	log.Printf("Register handler key %s handler %v", key, handler)
	k.handler[key] = handler
	return nil
}

func (k *kafkaConsumer) Start() {
	flag.StringVar(&kafkaBrokerUrl, "kafka-brokers", "localhost:19092,localhost:29092,localhost:39092", "Kafka brokers in comma separated value")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka topic.")

	flag.Parse()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	brokers := strings.Split(kafkaBrokerUrl, ",")

	config := kafka.ReaderConfig{
		Brokers:         brokers,
		Topic:           kafkaTopic,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	for {
		m, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		key := string(m.Key)

		if handler, ok := k.handler[key]; ok {
			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))

			event := handler.NewEvent()
			if err := json.Unmarshal(m.Value, event); err != nil {
				log.Printf("unmarshal error for key %s: %v", key, err)
				return
			}

			if err := handler.Handle(event); err != nil {
				log.Printf("handle error for key %s: %v", key, err)
			}
		}
	}
}
