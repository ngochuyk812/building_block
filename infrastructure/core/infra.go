package infrastructurecore

import (
	"strings"

	"github.com/ngochuyk812/building_block/infrastructure/cache"
	"github.com/ngochuyk812/building_block/infrastructure/eventbus"
	"github.com/ngochuyk812/building_block/infrastructure/eventbus/kafka"
)

type IInfra interface {
	InjectCache(connectString, pass string) error
	InjectEventbus(brokers, topic string) error
	GetCache() cache.ICache
	GetEventbus() eventbus.Producer
}
type Infra struct {
	cache    cache.ICache
	eventbus eventbus.Producer
}

var _ IInfra = (*Infra)(nil)

func NewInfra() IInfra {
	return &Infra{}
}

func (infra *Infra) GetEventbus() eventbus.Producer {
	return infra.eventbus
}

func (infra *Infra) InjectEventbus(brokers, topic string) error {
	if infra.eventbus != nil {
		return nil
	}
	eventbus, err := kafka.NewProceduer(strings.Split(brokers, ","), topic)
	if err != nil {
		panic("Cannot connect eventbus")
	}
	infra.eventbus = eventbus
	return nil
}

func (infra *Infra) InjectCache(connectString, pass string) error {
	if infra.cache != nil {
		return nil
	}
	cache, err := cache.NewRedisCache(connectString, pass)
	if err != nil {
		panic("Cannot connect redis")
	}
	infra.cache = cache
	return nil
}

func (infra *Infra) GetCache() cache.ICache {
	return infra.cache
}
