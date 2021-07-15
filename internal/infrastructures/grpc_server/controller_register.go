package grpc_server

import (
	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"google.golang.org/grpc"
)

type (
	controllerRegister struct {
		ticTacToeBattleController tictactoe_battle.TicTacToeBattleServiceServer
	}
)

func NewControllerRegister(controller tictactoe_battle.TicTacToeBattleServiceServer) ControllerRegister {
	return &controllerRegister{
		ticTacToeBattleController: controller,
	}
}

func (cr *controllerRegister) Register(grpcServer grpc.ServiceRegistrar) {
	tictactoe_battle.RegisterTicTacToeBattleServiceServer(grpcServer, cr.ticTacToeBattleController)
}
