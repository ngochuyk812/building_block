package mediator

type Mediator struct {
	handlers map[string]interface{}
}

func NewMediator() *Mediator {
	return &Mediator{
		handlers: make(map[string]interface{}),
	}
}

func (m *Mediator) AddHandler(key string, handler interface{}) {
	m.handlers[key] = handler
}

func (m *Mediator) GetHandler(key string) (interface{}, bool) {
	handler, ok := m.handlers[key]
	return handler, ok
}
