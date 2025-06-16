package main

import (
	"log/slog"
	"service-demo/internal/config"
	"service-demo/internal/database/pgsql"
	"service-demo/internal/database/redis"
	"service-demo/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		return
	}

	err = pgsql.Init(cfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer pgsql.CloseDB()

	err = redis.Init(cfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer redis.CloseDB()

	server.Start(cfg)
}
