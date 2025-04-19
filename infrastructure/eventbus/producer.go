package eventbus

import (
	"context"
)

type IntegrationEvent interface {
	Key() string
}

type Producer interface {
	Publish(ctx context.Context, event IntegrationEvent) error
}
