package ports

import (
	"context"

	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
)

type (
	BattleListener interface {
		Listen(ctx context.Context) (*tictactoe_battle.BattleSituation, error)
	}
)
