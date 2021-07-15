package domains

import (
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/battle"
)

type (
	Factory interface {
		BattleRule() battle.Rule
	}

	factory struct {
		battleRule battle.Rule
	}
)

func NewFactory() Factory {
	return &factory{
		battleRule: battle.NewRule(),
	}
}

func (f factory) BattleRule() battle.Rule {
	return f.battleRule
}
