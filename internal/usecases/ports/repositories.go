//go:generate mockgen -source=$GOFILE -destination=../../tests/mocks/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package ports

import (
	"context"

	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/room"
)

type (
	LoginRepository interface {
		FindByID(ctx context.Context, loginID string) (*tictactoe_battle.Login, error)
		NewLogin(ctx context.Context, loginID string) (*tictactoe_battle.Login, error)
		ReLogin(ctx context.Context, login *tictactoe_battle.Login) error
		Logout(ctx context.Context, loginID string) error
	}

	BattleRepository interface {
		Create(ctx context.Context, battle *battle.Battle) (room.ID, error)
		Update(ctx context.Context, battle *battle.Battle) error
		Enter(ctx context.Context, roomID room.ID, loginID string) error
		Leave(ctx context.Context, roomID room.ID, loginID string) error
		ReadStreamLatest(ctx context.Context, roomID room.ID) (string, *battle.Battle, error)
		ReadStream(ctx context.Context, roomID room.ID, messageID string) (string, *battle.Battle, error)
		ListMembers(ctx context.Context, roomID room.ID) ([]string, error)
		IsExistsInRoom(ctx context.Context, roomID room.ID, loginID string) (bool, error)
		Delete(ctx context.Context, roomID room.ID) error
	}
)
