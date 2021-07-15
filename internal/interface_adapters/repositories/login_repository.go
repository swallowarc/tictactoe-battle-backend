package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/exceptions"
	"github.com/swallowarc/tictactoe_battle_backend/internal/interface_adapters/gateways"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/ports"
	"golang.org/x/xerrors"
)

const (
	loginKeyPrefix = "tictactoe_battle_login"
	loginTimeout   = time.Hour
)

type (
	loginRepository struct {
		memDBCli gateways.MemDBClient
	}
)

func NewLoginRepository(gwFactory gateways.Factory) ports.LoginRepository {
	return &loginRepository{
		memDBCli: gwFactory.MemDBClient(),
	}
}

func (r *loginRepository) FindByID(ctx context.Context, loginID string) (*tictactoe_battle.Login, error) {
	sessionID, err := r.memDBCli.Get(ctx, loginKey(loginID))
	if err != nil {
		return nil, xerrors.Errorf("failed to memdb get: %w", err)
	}
	return &tictactoe_battle.Login{
		LoginId:   loginID,
		SessionId: sessionID,
	}, nil
}

func (r *loginRepository) NewLogin(ctx context.Context, loginID string) (*tictactoe_battle.Login, error) {
	sessionID := uuid.NewString()
	if err := r.memDBCli.SetNX(ctx, loginKey(loginID), sessionID, loginTimeout); err != nil {
		return nil, xerrors.Errorf("failed to SetNX: %w", err)
	}
	return &tictactoe_battle.Login{
		LoginId:   loginID,
		SessionId: sessionID,
	}, nil
}

func (r *loginRepository) ReLogin(ctx context.Context, login *tictactoe_battle.Login) error {
	if err := r.memDBCli.SetNX(ctx, loginKey(login.LoginId), login.SessionId, loginTimeout); err != nil {
		return xerrors.Errorf("failed to SetNX: %w", err)
	}
	return nil
}

func (r *loginRepository) Logout(ctx context.Context, loginID string) error {
	err := r.memDBCli.Del(ctx, loginKey(loginID))
	if exceptions.IsNotFoundError(err) {
		return nil
	}
	if err != nil {
		return xerrors.Errorf("failed to Del: %w", err)
	}
	return nil
}

func loginKey(loginID string) string {
	return fmt.Sprintf("%s:%s", loginKeyPrefix, loginID)
}
