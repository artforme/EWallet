package main

import (
	"EWallet/internal/config"
	"EWallet/internal/loggers"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup logger depending on env
	log := loggers.SetupLogger(cfg.Env)
}
