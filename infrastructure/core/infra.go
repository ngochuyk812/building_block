package infrastructurecore

import (
	"strings"

	"github.com/ngochuyk812/building_block/infrastructure/cache"
	"github.com/ngochuyk812/building_block/infrastructure/eventbus"
	"github.com/ngochuyk812/building_block/infrastructure/eventbus/kafka"
	"github.com/ngochuyk812/building_block/pkg/config"
	"go.uber.org/zap"
)

type IInfra interface {
	GetLogger() *zap.Logger
	InjectCache(connectString, pass string) error
	InjectEventbus(brokers, topic string) error
	GetCache() cache.ICache
	GetConfig() *config.ConfigApp
	GetEventbus() eventbus.Producer
}
type Infra struct {
	config   *config.ConfigApp
	logger   *zap.Logger
	cache    cache.ICache
	eventbus eventbus.Producer
}

var _ IInfra = (*Infra)(nil)

func NewInfra(config *config.ConfigApp) IInfra {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("cannnot install logger: " + err.Error())
	}
	infra := &Infra{
		config: config,
		logger: logger,
	}
	return infra
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

func (infra *Infra) GetLogger() *zap.Logger {
	return infra.logger
}

func (infra *Infra) GetCache() cache.ICache {
	return infra.cache
}

func (infra *Infra) GetConfig() *config.ConfigApp {
	return infra.config
}
