package bus_core

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ngochuyk812/building_block/infrastructure/helpers"
	"github.com/ngochuyk812/building_block/pkg/mediator"
)

type IHandler[RQ any, RES any] interface {
	Handle(ctx context.Context, cmd RQ) (RES, error)
}

func RegisterHandler[RQ any, RES any](m *mediator.Mediator, cmd RQ, handler IHandler[RQ, RES]) {
	key := fmt.Sprintf("%T", cmd)
	m.AddHandler(key, handler)
}

func Send[C any, R any](m *mediator.Mediator, ctx context.Context, cmd C) (res R, err error) {

	key := fmt.Sprintf("%T", cmd)

	authContext, ok := helpers.AuthContext(ctx)
	if ok == false {
		log.Fatalf("cannot get authcontext")
	}
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		log.Printf(
			"[%s] siteId=%s username=%s request=%+v duration=%dms res=%+v err=%+v",
			key,
			authContext.IdSite,
			authContext.UserName,
			cmd,
			duration,
			res,
			err,
		)
	}()
	handler, ok := m.GetHandler(key)
	if !ok {
		// res.AddError(400, fmt.Sprintf("no handler registered for type %T", cmd))
		return res, fmt.Errorf("no handler registered for type %T", cmd)
	}

	h, ok := handler.(IHandler[C, R])
	if !ok {
		// res.AddError(400, fmt.Sprintf("handler type mismatch for type %T", cmd))
		return res, fmt.Errorf("handler type mismatch for type %T", cmd)
	}
	res, err = h.Handle(ctx, cmd)

	return res, err
}
