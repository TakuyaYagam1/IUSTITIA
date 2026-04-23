package main

import (
	"log"

	logkit "github.com/wahrwelt-kit/go-logkit"

	"github.com/TakuyaYagam1/iustitia/config"
	"github.com/TakuyaYagam1/iustitia/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config initialization failed: %v", err)
	}

	l, err := logkit.New(
		logkit.WithLevel(logkit.InfoLevel),
		logkit.WithOutput(logkit.ConsoleOutput),
		logkit.WithServiceName("iustitia-backend"),
	)
	if err != nil {
		log.Fatalf("logger initialization failed: %v", err)
	}

	l.Info("configuration loaded", logkit.Fields{"config": cfg.String()})

	app.Run(cfg, l)
}
