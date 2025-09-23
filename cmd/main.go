package main

import (
	"api-gateway-sql/config"
	_ "api-gateway-sql/docs"
	"api-gateway-sql/internal/app"
	"api-gateway-sql/pkg/logger"

	"github.com/go-playground/validator/v10"
)

// @title API GATEWAY SQL
// @description API used for executing SQL QUERY
// @contact.name API Support
// @contact.email ngaswilly77@gmail.com
// @license.name MIT
// @license.url https://github.com/willbrid/api_gateway_sql/blob/main/LICENSE
// @BasePath /
// @securityDefinitions.basic BasicAuth
func main() {
	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

	configFlag, err := config.LoadConfigFlag(validate)
	if err != nil {
		logger.Fatal("failed to load configuration flags: %v", err.Error())
	}

	viperInstance, err := config.ReadConfigFile(configFlag.ConfigFile)
	if err != nil {
		logger.Fatal("failed to read configuration file: %v", err.Error())
	}

	configLoaded, err := config.LoadConfig(viperInstance, validate)
	if err != nil {
		logger.Fatal("failed to load configuration file: %v", err.Error())
	}

	logger.Info("configuration file '%s' was loaded successfully", configFlag.ConfigFile)

	app.Run(configLoaded, configFlag)
}
