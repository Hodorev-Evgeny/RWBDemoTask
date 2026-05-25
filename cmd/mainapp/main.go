package main

import (
	core_logger "RWBDwmoTask/internal/core/logger"
	core_server "RWBDwmoTask/internal/core/transport/server"
	"context"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	serverConfig := core_server.MustNewConfigServer()
	time.Local = serverConfig.TimeZone

	loggerConfig := core_logger.MustNewConfig()
	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Debug("starting server")

	apiVersionRouters := core_server.NewAPIVersionRouter(core_server.ApiVersion1)

	server := core_server.NewServer(serverConfig, logger)

	server.ResisterApiVersionRouter(apiVersionRouters)
	server.AddFrond()

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start server", zap.Error(err))
	}
}
