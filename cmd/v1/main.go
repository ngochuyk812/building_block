package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"

	infrastructurecore "github.com/ngochuyk812/building_block/infrastructure/core"
	"github.com/ngochuyk812/building_block/pkg/config"
)

func main() {

	config := config.NewConfigEnv()

	infa := infrastructurecore.NewInfra(config)
	infa.InjectCache(config.RedisConnect, config.RedisPass)

}
