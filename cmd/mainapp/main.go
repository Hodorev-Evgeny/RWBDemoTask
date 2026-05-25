package main

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	core_logger "RWBDwmoTask/internal/core/logger"
	core_nats "RWBDwmoTask/internal/core/repository/nats"
	core_redis "RWBDwmoTask/internal/core/repository/redis"
	core_server "RWBDwmoTask/internal/core/transport/server"
	feature_repository_toplist "RWBDwmoTask/internal/features/toplist/repository"
	feature_service_toplist "RWBDwmoTask/internal/features/toplist/service"
	feature_transport_toplist "RWBDwmoTask/internal/features/toplist/transport"
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
	natsConfig := core_nats.MustNewNatsConfig()
	natsClient, err := core_nats.New(ctx, natsConfig)
	if err != nil {
		panic(err)
	}

	redisConfig := core_redis.MustGetRedisConfig()
	redisClient := core_redis.CreateRedisClientMust(redisConfig)

	storeg := core_domain.NewStorage(
		5*time.Minute,
		5*time.Second,
		redisClient,
	)
	go storeg.Run()

	listrepositoy := feature_repository_toplist.NewNatsRepository(storeg)
	listservise := feature_service_toplist.NewNatsService(listrepositoy)
	listtransport := feature_transport_toplist.NewTransportTopList(listservise)
	listRoute := listtransport.Router()

	apiVersionRouters := core_server.NewAPIVersionRouter(core_server.ApiVersion1)

	server := core_server.NewServer(serverConfig, logger, storeg)

	server.ResisterApiVersionRouter(apiVersionRouters)
	server.RegisterRoutes(listRoute...)

	go func() {
		if err := server.ReadEvents(ctx, natsClient.JS); err != nil {
			logger.Error("read nats events error", zap.Error(err))
		}
	}()

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start server", zap.Error(err))
	}
}
