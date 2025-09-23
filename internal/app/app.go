package app

import (
	"api-gateway-sql/config"
	"api-gateway-sql/internal/delivery"
	"api-gateway-sql/internal/repository"
	"api-gateway-sql/internal/usecase"
	"api-gateway-sql/pkg/database"
	"api-gateway-sql/pkg/httpserver"
	"api-gateway-sql/pkg/logger"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag) {
	sqliteAppDatabase := database.NewSqliteAppDatabase(cfgfile.ApiGatewaySQL.Sqlitedb)
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

	logger.LogInfo(logInfoServer, cfgflag.ListenPort)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.LogInfo("app server - run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		logger.LogError("app server error: %v", err.Error())
	}

	if err := httpServer.Stop(); err != nil {
		logger.LogError("app server - stop - error: %v", err)
	}
}
