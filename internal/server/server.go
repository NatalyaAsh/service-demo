package server

import (
	"log/slog"
	"net/http"
	"os"

	"service-demo/internal/api"
	"service-demo/internal/config"
)

func Start(cfg *config.Config) {

	mux := http.NewServeMux()

	api.Init(mux, cfg)

	slog.Info("Start server")
	err := http.ListenAndServe(":"+cfg.HTTP.Port, mux)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

}
