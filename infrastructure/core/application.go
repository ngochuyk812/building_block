package infrastructurecore

import (
	"net/http"

	"go.uber.org/zap"
)

type ServeMux struct {
	port   string
	Mux    *http.ServeMux
	logger *zap.Logger
}

func NewServe(port string, logger *zap.Logger) *ServeMux {
	serve := &ServeMux{
		port:   port,
		Mux:    http.NewServeMux(),
		logger: logger,
	}

	return serve
}

func (app *ServeMux) Run() {
	app.logger.Info("server start", zap.String("port", app.port))
	err := http.ListenAndServe(app.port, app.Mux)
	if err != nil {
		app.logger.Error("server start error", zap.String("port", app.port), zap.Error(err))
		panic(err)
	}
}
