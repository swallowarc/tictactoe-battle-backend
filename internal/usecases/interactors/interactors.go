//go:generate mockgen -source=$GOFILE -destination=../../tests/mocks/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package interactors

import (
	"context"

	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/room"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/ports"
)

type (
	LoginInteractor interface {
		Login(ctx context.Context, login *tictactoe_battle.Login) (*tictactoe_battle.Login, error)
		Logout(ctx context.Context, login *tictactoe_battle.Login) error
	}

	BattleInteractor interface {
		Create(ctx context.Context) (room.ID, error)
		CanEnter(ctx context.Context, roomID room.ID, loginID string) (bool, error)
		Enter(ctx context.Context, roomID room.ID, loginID string) (ports.BattleListener, error)
		Declaration(ctx context.Context, roomID room.ID, loginID string) error
		Leave(ctx context.Context, roomID room.ID, loginID string) error
		Attack(ctx context.Context, roomID room.ID, player tictactoe_battle.Player, position tictactoe_battle.Position, pieceSize tictactoe_battle.Piece) error
		Pick(ctx context.Context, roomID room.ID, player tictactoe_battle.Player, position tictactoe_battle.Position, pieceSize tictactoe_battle.Piece) error
		Reset(ctx context.Context, roomID room.ID) error
	}
)
