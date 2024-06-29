package main

import (
	"github.com/chernyshevuser/gopfermart.git/tools/closer"
	"github.com/chernyshevuser/gopfermart.git/tools/config"
	"github.com/chernyshevuser/gopfermart.git/tools/logger"
)

func main() {
	logger := logger.New()
	defer logger.Sync()

	config.SetupConfig(logger)

	go closer.GracefulShutdown(
		make(chan struct{}, 1),
	)
}
