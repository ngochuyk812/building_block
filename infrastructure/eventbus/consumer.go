package eventbus

import (
	"context"
)

type IntegrationEventHandler interface {
	NewEvent() IntegrationEvent
	Handle(ctx context.Context, event IntegrationEvent) error
}

type Consumer interface {
	RegisterHandler(handler IntegrationEventHandler) (err error)
	Run()
}
