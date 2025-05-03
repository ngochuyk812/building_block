package mockup

import (
	"context"
	"fmt"
	"testing"

	infrastructurecore "github.com/ngochuyk812/building_block/infrastructure/core"
	"github.com/ngochuyk812/building_block/infrastructure/eventbus"
	"github.com/ngochuyk812/building_block/infrastructure/eventbus/kafka"
)

func TestNewConsumer(t *testing.T) {
	consumer, err := kafka.NewConsumer("localhost:9092", "topic-test", "group-test")
	if err != nil {
		t.Error(err)
		return
	}
	consumer.RegisterHandler(NewEventSendMailHandler())

	consumer.Run()
}

func NewEventSendMailHandler() eventbus.IntegrationEventHandler {
	return &eventSendMailHandler{
		infra: infrastructurecore.NewInfra(),
	}
}

type EventSendMail struct {
	Content string
	Title   string
}

func (e *EventSendMail) Key() string {
	return "test"
}

type eventSendMailHandler struct {
	infra infrastructurecore.IInfra
}

func (e *eventSendMailHandler) NewEvent() eventbus.IntegrationEvent {
	return &EventSendMail{}
}
func (e *eventSendMailHandler) Handle(ctx context.Context, event eventbus.IntegrationEvent) error {
	sendMailEvent, oke := event.(*EventSendMail)
	if oke == false {
	}
	fmt.Println(sendMailEvent.Content)
	return nil
}
