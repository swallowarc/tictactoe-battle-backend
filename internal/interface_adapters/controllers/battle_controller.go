package controllers

import (
	"context"
	"time"

	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/exceptions"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/loggers"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/room"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/listener"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func (c *ticTacToeBattleController) CreateRoom(ctx context.Context, _ *tictactoe_battle.CreateRoomRequest) (*tictactoe_battle.CreateRoomResponse, error) {
	roomID, err := c.battleInteractor.Create(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to Create: %w", err)
	}

	return &tictactoe_battle.CreateRoomResponse{
		RoomId: roomID.String(),
	}, nil
}

func (c *ticTacToeBattleController) CanEnterRoom(ctx context.Context, req *tictactoe_battle.CanEnterRoomRequest) (*tictactoe_battle.CanEnterRoomResponse, error) {
	can, err := c.battleInteractor.CanEnter(ctx, room.ID(req.RoomId), req.LoginId)
	if err != nil {
		return nil, xerrors.Errorf("failed to CanEnter: %w", err)
	}

	return &tictactoe_battle.CanEnterRoomResponse{CanEnterRoom: can}, nil
}

func (c *ticTacToeBattleController) EnterRoom(request *tictactoe_battle.EnterRoomRequest, stream tictactoe_battle.TicTacToeBattleService_EnterRoomServer) error {
	ctx := loggers.LoggerToContext(stream.Context(), c.logger)
	lsnr, err := c.battleInteractor.Enter(ctx, room.ID(request.RoomId), request.LoginId)
	if err != nil {
		return xerrors.Errorf("failed to Enter: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			bt, err := lsnr.Listen(ctx)
			if err != nil {
				if exceptions.IsNotFoundError(err) {
					time.Sleep(time.Second)
					continue
				}
				if xerrors.As(err, &listener.LeftError) {
					loggers.Logger(ctx).Info("already left the room")
					return nil
				}
				if xerrors.As(err, &context.Canceled) {
					loggers.Logger(ctx).Info("context canceled")
					return nil
				}

				loggers.Logger(ctx).Error("failed to Listen", zap.Error(err))
				return xerrors.Errorf("failed to Listen: %w", err)
			}
			if err := stream.Send(bt); err != nil {
				if xerrors.As(err, &context.Canceled) {
					loggers.Logger(ctx).Debug("client context canceled")
					return nil
				}

				return xerrors.Errorf("failed to Send: %w", err)
			}
		}
	}
}

func (c *ticTacToeBattleController) Declaration(ctx context.Context, req *tictactoe_battle.DeclarationRequest) (*tictactoe_battle.NoBody, error) {
	if err := c.battleInteractor.Declaration(ctx, room.ID(req.RoomId), req.LoginId); err != nil {
		return nil, xerrors.Errorf("failed to Declaration: %w", err)
	}

	return &tictactoe_battle.NoBody{}, nil
}

func (c *ticTacToeBattleController) LeaveRoom(ctx context.Context, req *tictactoe_battle.LeaveRoomRequest) (*tictactoe_battle.NoBody, error) {
	if err := c.battleInteractor.Leave(ctx, room.ID(req.RoomId), req.LoginId); err != nil {
		return nil, xerrors.Errorf("failed to Leave: %w", err)
	}

	return &tictactoe_battle.NoBody{}, nil
}

func (c *ticTacToeBattleController) Attack(ctx context.Context, req *tictactoe_battle.AttackRequest) (*tictactoe_battle.NoBody, error) {
	if err := c.battleInteractor.Attack(ctx, room.ID(req.RoomId), req.Player, req.Position, req.Piece); err != nil {
		return nil, xerrors.Errorf("failed to Attack: %w", err)
	}

	return &tictactoe_battle.NoBody{}, nil
}

func (c *ticTacToeBattleController) Pick(ctx context.Context, req *tictactoe_battle.PickRequest) (*tictactoe_battle.NoBody, error) {
	if err := c.battleInteractor.Pick(ctx, room.ID(req.RoomId), req.Player, req.Position, req.Piece); err != nil {
		return nil, xerrors.Errorf("failed to Pick: %w", err)
	}

	return &tictactoe_battle.NoBody{}, nil
}

func (c *ticTacToeBattleController) ResetBattle(ctx context.Context, req *tictactoe_battle.ResetBattleRequest) (*tictactoe_battle.NoBody, error) {
	if err := c.battleInteractor.Reset(ctx, room.ID(req.RoomId)); err != nil {
		return nil, xerrors.Errorf("failed to Reset: %w", err)
	}

	return &tictactoe_battle.NoBody{}, nil
}
