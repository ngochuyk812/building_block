package main

import (
	infrastructurecore "building_block/infrastructure/core"
	"building_block/infrastructure/databases"
	"building_block/internal/interfaces/v1/connectrpc"
	"building_block/pkg/config"
	"fmt"
	"os"
	"os/signal"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Test2 struct {
	Key   string        `redis:key;json:key`
	Value time.Duration `redis:value;json:value`
}

func main() {
	policiesPath := &map[string][]string{
		"/greet.v1.GreetService/Greet": {"user"},
	}
	config := config.NewConfigEnv()
	config.PoliciesPath = policiesPath

	infa := infrastructurecore.NewInfra(config)
	infa.InjectSQL(databases.MYSQL)
	infa.InjectCache(config.RedisConnect, config.RedisPass)
	app1 := infrastructurecore.NewServe(":"+config.Port, infa.GetLogger())

	path, handler := connectrpc.NewGreetServer(infa)
	app1.Mux.Handle(path, handler)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go app1.Run()
	<-c
	fmt.Println("shutting down...")

}
