package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"

	infrastructurecore "github.com/ngochuyk812/building_block/infrastructure/core"
	"github.com/ngochuyk812/building_block/infrastructure/databases"
	"github.com/ngochuyk812/building_block/pkg/config"
)

func main() {
	policiesPath := &map[string][]string{
		"/greet.v1.GreetService/Greet": {"user"},
	}
	config := config.NewConfigEnv()
	config.PoliciesPath = policiesPath

	infa := infrastructurecore.NewInfra(config)
	infa.InjectSQL(databases.MYSQL)
	infa.InjectCache(config.RedisConnect, config.RedisPass)

}
