package infrastructurecore

import (
	"building_block/infrastructure/cache"
	"building_block/infrastructure/databases"
	"building_block/pkg/config"
	"database/sql"

	"go.uber.org/zap"
)

type IInfra interface {
	GetLogger() *zap.Logger
	GetDb() *sql.DB
	InjectSQL(dbType databases.DatabaseType) error
	InjectCache(connectString, pass string) error
	GetCache() cache.ICache
	GetConfig() *config.ConfigApp
}
type Infra struct {
	config *config.ConfigApp
	logger *zap.Logger
	dbSql  *sql.DB
	cache  cache.ICache
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
	if infra.dbSql != nil {
		return nil
	}

	db, err := databases.NewConnectSql(dbType, infra.config.DbName, infra.config.DbConnect, infra.logger)
	if err != nil {
		panic("Error connect sql: " + err.Error())
	}
	infra.dbSql = db
	if infra.dbSql == nil {
		panic("Error conenct sql")
	}
	return nil
}
func (infra *Infra) GetDb() *sql.DB {
	return infra.dbSql
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
