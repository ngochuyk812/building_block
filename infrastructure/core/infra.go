package infrastructurecore

import (
	"github.com/ngochuyk812/building_block/infrastructure/cache"
	"github.com/ngochuyk812/building_block/infrastructure/databases"
	"github.com/ngochuyk812/building_block/pkg/config"
	"github.com/ngochuyk812/building_block/pkg/mediator"
	"go.uber.org/zap"
)

type IInfra interface {
	GetLogger() *zap.Logger
	GetDatabase() databases.IDatabase
	InjectSQL(dbType databases.DatabaseType) error
	InjectCache(connectString, pass string) error
	GetCache() cache.ICache
	GetConfig() *config.ConfigApp
	GetMediator() *mediator.Mediator
}
type Infra struct {
	config   *config.ConfigApp
	logger   *zap.Logger
	cache    cache.ICache
	database databases.IDatabase
	mediator *mediator.Mediator
}

var _ IInfra = (*Infra)(nil)

func NewInfra(config *config.ConfigApp) IInfra {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("cannnot install logger: " + err.Error())
	}
	infra := &Infra{
		config:   config,
		logger:   logger,
		mediator: mediator.NewMediator(),
	}
	return infra
}
func (infra *Infra) GetDatabase() databases.IDatabase {
	return infra.database
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

func (infra *Infra) InjectSQL(dbType databases.DatabaseType) error {
	if infra.database != nil {
		return nil
	}
	database, err := databases.NewDatabases(dbType, infra.config.DbConnectRead, infra.config.DbConnect, infra.config.DbName, infra.logger)
	if err != nil {
		panic("Error connect sql: " + err.Error())
	}
	infra.database = database
	if infra.database == nil {
		panic("Error conenct sql")
	}
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

func (infra *Infra) GetMediator() *mediator.Mediator {
	return infra.mediator
}
