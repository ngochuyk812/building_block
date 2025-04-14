package bus_core

import (
	"context"
	"fmt"

	"github.com/ngochuyk812/building_block/pkg/mediator"
)

type IHandler[RQ any, RES any] interface {
	Handle(ctx context.Context, cmd RQ) (RES, error)
}

func RegisterHandler[RQ any, RES any](m *mediator.Mediator, cmd RQ, handler IHandler[RQ, RES]) {
	key := fmt.Sprintf("%T", cmd)
	m.AddHandler(key, handler)
}

func Send[C any, R any](m *mediator.Mediator, ctx context.Context, cmd C) (R, error) {
	key := fmt.Sprintf("%T", cmd)

	handler, ok := m.GetHandler(key)
	if !ok {
		var zero R
		return zero, fmt.Errorf("no handler registered for type %T", cmd)
	}

	h, ok := handler.(IHandler[C, R])
	if !ok {
		var zero R
		return zero, fmt.Errorf("handler type mismatch for type %T", cmd)
	}

	return h.Handle(ctx, cmd)
}
