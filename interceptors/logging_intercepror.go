package interceptors

import (
	"context"
	"encoding/json"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

func NewLoggingInterceptor(logger *zap.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (res connect.AnyResponse, err error) {
			start := time.Now()
			response, errService := next(ctx, req)

			duration := time.Since(start).Milliseconds()
			if errService != nil {
				logger.Error(
					req.Spec().Procedure,
					zap.Error(errService),
					zap.Int64("took", duration),
				)
			} else {
				respBytes, _ := json.Marshal(response.Any())

				logger.Info(
					req.Spec().Procedure,
					zap.ByteString("response", respBytes),
					zap.Int64("took", duration),
				)
			}
			return response, errService
		}
	}
}
