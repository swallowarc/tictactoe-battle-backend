package controllers

import (
	"context"

	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/loggers"
	"golang.org/x/xerrors"
)

func (c *ticTacToeBattleController) Login(ctx context.Context, request *tictactoe_battle.LoginRequest) (*tictactoe_battle.LoginResponse, error) {
	login := request.Login
	loggers.With(ctx, loggers.Map{
		"login_id":   login.LoginId,
		"session_id": login.SessionId,
	})

	newLogin, err := c.loginInteractor.Login(ctx, login)
	if err != nil {
		return nil, xerrors.Errorf("failed to login: %w", err)
	}

	return &tictactoe_battle.LoginResponse{
		Login: newLogin,
	}, nil
}

func (c *ticTacToeBattleController) Logout(ctx context.Context, request *tictactoe_battle.LogoutRequest) (*tictactoe_battle.NoBody, error) {
	login := request.Login
	loggers.With(ctx, loggers.Map{
		"login_id":   login.LoginId,
		"session_id": login.SessionId,
	})
	if err := c.loginInteractor.Logout(ctx, login); err != nil {
		return nil, xerrors.Errorf("failed to logout: %w", err)
	}
	return &tictactoe_battle.NoBody{}, nil
}
