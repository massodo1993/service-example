package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/massodo1993/service-example/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	MongoConfig   MongoConfig
	IamGRPC       IamGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryGRPCConfig, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongoConfig, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	iamGRPCConfig, err := env.NewIamGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		InventoryGRPC: inventoryGRPCConfig,
		MongoConfig:   mongoConfig,
		IamGRPC:       iamGRPCConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
