package main

import (
	"github.com/willbrid/api-gateway-sql/config"
	_ "github.com/willbrid/api-gateway-sql/docs"
	"github.com/willbrid/api-gateway-sql/internal/app"
	"github.com/willbrid/api-gateway-sql/pkg/logging"

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
	validate := validator.New(validator.WithRequiredStructEnabled())
	logger := logging.InitLogger()

	configFlag, err := config.LoadConfigFlag(validate)
	if err != nil {
		logger.Error().Err(err).Msg("failed to load configuration flags")
		return
	}

	viperInstance, err := config.ReadConfigFile(configFlag.ConfigFile)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read configuration file")
		return
	}

	configLoaded, err := config.LoadConfig(viperInstance, validate)
	if err != nil {
		logger.Error().Err(err).Msg("failed to load configuration file")
		return
	}

	logger.Info().Str("config_file", configFlag.ConfigFile).Msg("configuration file was loaded successfully")
	app.Run(configLoaded, configFlag, logger)
}
