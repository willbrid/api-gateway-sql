package app

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/delivery"
	"github.com/willbrid/api-gateway-sql/internal/delivery/middleware"
	"github.com/willbrid/api-gateway-sql/internal/repository"
	"github.com/willbrid/api-gateway-sql/internal/usecase"
	"github.com/willbrid/api-gateway-sql/pkg/csvstream"
	"github.com/willbrid/api-gateway-sql/pkg/database"
	"github.com/willbrid/api-gateway-sql/pkg/httpserver"
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag, loggerInstance logger.ILogger) {
	sqliteAppDatabase, err := database.NewSqliteAppDatabase(cfgfile.ApiGatewaySQL.Sqlitedb)
	if err != nil {
		loggerInstance.Error("app server database error: %v", err.Error())
		return
	}
	MigrateAppDatabase(sqliteAppDatabase.Db, loggerInstance)

	repos := repository.NewRepositories(sqliteAppDatabase.Db)
	csvstream := csvstream.NewCSVStream(loggerInstance)
	usecases := usecase.NewUsecases(usecase.Deps{
		Repos:      repos,
		Config:     cfgfile,
		ICSVStream: csvstream,
		ILogger:    loggerInstance,
	})

	httpServer := httpserver.NewServer(
		fmt.Sprint(":"+fmt.Sprint(cfgflag.ListenPort)),
		cfgflag.EnableHttps,
		cfgflag.CertFile,
		cfgflag.KeyFile,
	)
	authMiddleware := middleware.NewAuthMiddleware(loggerInstance)
	handlers := delivery.NewHandler(usecases, httpServer, authMiddleware, loggerInstance)
	handlers.InitRouter(cfgfile, cfgflag)
	httpServer.Start()

	scheme := map[bool]string{true: "https", false: "http"}[cfgflag.EnableHttps]
	loggerInstance.Info("app server is listening on port %v using %s", cfgflag.ListenPort, scheme)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		loggerInstance.Info("app server - run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		loggerInstance.Error("app server error: %v", err.Error())
	}

	if appDbCnx, err := sqliteAppDatabase.Db.DB(); err == nil {
		_ = appDbCnx.Close()
	}

	if err := httpServer.Stop(); err != nil {
		loggerInstance.Error("app server - stop - error: %v", err)
	}
}
