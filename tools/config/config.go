package config

import (
	"flag"

	"github.com/chernyshevuser/gopfermart.git/tools/logger"
)

type configKey string

const (
	RunAddrEnv           = configKey("RUN_ADDRESS")
	DatabaseUriEnv       = configKey("DATABASE_URI")
	AccrualSystemAddrEnv = configKey("ACCRUAL_SYSTEM_ADDRESS")
)

var (
	RunAddr           string
	DatabaseUri       string
	AccrualSystemAddr string
)

func SetupConfig(logger logger.Logger) {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "runAddr")
	flag.StringVar(&DatabaseUri, "d", "", "dbUri")
	flag.StringVar(&AccrualSystemAddr, "r", "", "accrual system addr")

	flag.Parse()

	runAddr, err := GetConfigString(RunAddrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		RunAddr = runAddr
	}

	databaseUri, err := GetConfigString(DatabaseUriEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		DatabaseUri = databaseUri
	}

	accrualSystemAddr, err := GetConfigString(AccrualSystemAddrEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		AccrualSystemAddr = accrualSystemAddr
	}

	logger.Infow(
		"config",
		"runAddr", RunAddr,
		"dbUri", DatabaseUri,
		"accrual system addr", AccrualSystemAddr,
	)
}
