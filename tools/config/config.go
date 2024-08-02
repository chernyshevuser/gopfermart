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
	CryptoKeyEnv         = configKey("CRYPTO_KEY")
	JwtSecretKeyEnv      = configKey("JWT_KEY")
)

var (
	RunAddr           string
	DatabaseUri       string
	AccrualSystemAddr string
	CryptoKey         string
	JwtSecretKey      string
)

func SetupConfig(logger logger.Logger) {
	flag.StringVar(&RunAddr, "a", "localhost:8080", "runAddr")
	flag.StringVar(&DatabaseUri, "d", "", "dbUri")
	flag.StringVar(&AccrualSystemAddr, "r", "", "accrual system addr")
	flag.StringVar(&CryptoKey, "k", "examplekey123456", "crypto key")
	flag.StringVar(&JwtSecretKey, "j", "examplekey123456", "jwt crypto key")

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

	cryptoKey, err := GetConfigString(CryptoKeyEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		CryptoKey = cryptoKey
	}

	jwtSecretKey, err := GetConfigString(JwtSecretKeyEnv)
	if err != nil {
		logger.Errorw(
			"can't get env",
			"msg", err,
		)
	} else {
		JwtSecretKey = jwtSecretKey
	}

	//TODO printing envs lol
	logger.Infow(
		"config",
		"runAddr", RunAddr,
		"dbUri", DatabaseUri,
		"accrual system addr", AccrualSystemAddr,
		"crypto key", CryptoKey,
		"jwt crypto key", JwtSecretKey,
	)
}
