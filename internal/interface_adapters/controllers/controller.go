package controllers

import (
	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/interactors"
	"go.uber.org/zap"
)

type (
	ticTacToeBattleController struct {
		logger *zap.Logger
		tictactoe_battle.UnimplementedTicTacToeBattleServiceServer

		loginInteractor  interactors.LoginInteractor
		battleInteractor interactors.BattleInteractor
	}
)

func NewTicTacToeBattleController(logger *zap.Logger, iFactory interactors.Factory) tictactoe_battle.TicTacToeBattleServiceServer {
	return &ticTacToeBattleController{
		logger:           logger,
		loginInteractor:  iFactory.LoginInteractor(),
		battleInteractor: iFactory.BattleInteractor(),
	}
}
