package main

import (
	"context"

	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/loggers"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains"
	"github.com/swallowarc/tictactoe_battle_backend/internal/infrastructures"
	"github.com/swallowarc/tictactoe_battle_backend/internal/infrastructures/env"
	"github.com/swallowarc/tictactoe_battle_backend/internal/infrastructures/grpc_server"
	"github.com/swallowarc/tictactoe_battle_backend/internal/interface_adapters/controllers"
	"github.com/swallowarc/tictactoe_battle_backend/internal/interface_adapters/repositories"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/interactors"
	"go.uber.org/zap"
)

func setup() grpc_server.GRPCServer {
	zapLogger := loggers.NewZapLogger(env.Server.RunMode)

	// factories
	dFactory := domains.NewFactory()
	gwFactory := infrastructures.NewFactory()
	repoFactory := repositories.NewFactory(gwFactory)
	iFactory := interactors.NewFactory(dFactory, repoFactory)

	// interface_adapters
	controller := controllers.NewTicTacToeBattleController(zapLogger, iFactory)
	// grpc_service_register
	grpcServiceRegister := grpc_server.NewControllerRegister(controller)

	// initializer & closer
	init := func() {
		if err := gwFactory.MemDBClient().Ping(context.Background()); err != nil {
			zapLogger.Panic("failed to ping to redis", zap.Error(err))
		}
		zapLogger.Info("ping to redis was successful")
	}
	closer := func() {}

	grpcServer := grpc_server.NewGRPCServer(
		zapLogger,
		env.Server.PORT,
		env.Server.RunMode,
		grpcServiceRegister,
		init,
		closer,
	)

	return grpcServer
}
