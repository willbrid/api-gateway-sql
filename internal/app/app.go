package app

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/delivery"
	"github.com/willbrid/api-gateway-sql/internal/delivery/middleware"
	"github.com/willbrid/api-gateway-sql/internal/repository"
	"github.com/willbrid/api-gateway-sql/internal/usecase"
	"github.com/willbrid/api-gateway-sql/pkg/database"
	"github.com/willbrid/api-gateway-sql/pkg/httpserver"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag, logger zerolog.Logger) {
	sqliteAppDatabase, err := database.NewSqliteAppDatabase(cfgfile.ApiGatewaySQL.Sqlitedb)
	if err != nil {
		logger.Error().Err(err).Msg("failed to init database server")
		return
	}

	if err := MigrateAppDatabase(sqliteAppDatabase.Db); err != nil {
		logger.Error().Err(err).Msg("failed to init database server")
		return
	}

	repos := repository.NewRepositories(sqliteAppDatabase.Db, logger)
	usecases := usecase.NewUsecases(usecase.Deps{
		Repos:  repos,
		Config: cfgfile,
		Logger: logger,
	})

	httpServer := httpserver.NewServer(
		fmt.Sprint(":"+fmt.Sprint(cfgflag.ListenPort)),
		cfgflag.EnableHttps,
		cfgflag.CertFile,
		cfgflag.KeyFile,
	)
	authMiddleware := middleware.NewAuthMiddleware(logger)
	handlers := delivery.NewHandler(usecases, httpServer, authMiddleware, logger)
	handlers.InitRouter(cfgfile, cfgflag)
	httpServer.Start()

	scheme := map[bool]string{true: "https", false: "http"}[cfgflag.EnableHttps]
	logger.Info().Str("scheme", scheme).Int("port", cfgflag.ListenPort).Msg("app server starting")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info().Str("signal", s.String()).Msg("app server stopping")
	case err := <-httpServer.Notify():
		logger.Error().Err(err).Msg("app server stopping")
	}

	if appDbCnx, err := sqliteAppDatabase.Db.DB(); err == nil {
		_ = appDbCnx.Close()
	}

	if err := httpServer.Stop(); err != nil {
		logger.Error().Err(err).Msg("app server stopping")
	}
}
