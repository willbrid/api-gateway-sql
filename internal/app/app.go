package app

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/delivery"
	"github.com/willbrid/api-gateway-sql/internal/repository"
	"github.com/willbrid/api-gateway-sql/internal/usecase"
	"github.com/willbrid/api-gateway-sql/pkg/database"
	"github.com/willbrid/api-gateway-sql/pkg/httpserver"
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag) {
	sqliteAppDatabase, err := database.NewSqliteAppDatabase(cfgfile.ApiGatewaySQL.Sqlitedb)
	if err != nil {
		logger.Error("app server database error: %v", err.Error())
		return
	}
	MigrateAppDatabase(sqliteAppDatabase.Db)

	repos := repository.NewRepositories(sqliteAppDatabase.Db)
	usecases := usecase.NewUsecases(usecase.Deps{
		Repos:  repos,
		Config: cfgfile,
	})

	handlers := delivery.NewHandler(usecases)
	httpServer := httpserver.NewServer(
		fmt.Sprint(":"+fmt.Sprint(cfgflag.ListenPort)),
		cfgflag.EnableHttps,
		cfgflag.CertFile,
		cfgflag.KeyFile,
	)
	handlers.InitRouter(httpServer.Router, cfgfile, cfgflag)
	httpServer.Start()

	var logInfoServer string
	if cfgflag.EnableHttps {
		logInfoServer = "app server is listening on port %v using https"
	} else {
		logInfoServer = "app server is listening on port %v using http"
	}

	logger.Info(logInfoServer, cfgflag.ListenPort)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app server - run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		logger.Error("app server error: %v", err.Error())
	}

	if appDbCnx, err := sqliteAppDatabase.Db.DB(); err == nil {
		appDbCnx.Close()
	}

	if err := httpServer.Stop(); err != nil {
		logger.Error("app server - stop - error: %v", err)
	}
}
