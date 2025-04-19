package eventbus

type IntegrationEventHandler interface {
	NewEvent() IntegrationEvent
	Handle(event IntegrationEvent) error
}

type Consumer interface {
	RegisterHandler(handler IntegrationEventHandler) (err error)
}
