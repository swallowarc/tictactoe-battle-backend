package interactors

import (
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/ports"
)

type (
	Factory interface {
		LoginInteractor() LoginInteractor
		BattleInteractor() BattleInteractor
	}

	factory struct {
		loginInteractor  LoginInteractor
		battleInteractor BattleInteractor
	}
)

func NewFactory(dFactory domains.Factory, rFactory ports.RepositoriesFactory) Factory {
	return &factory{
		loginInteractor:  NewLoginInteractor(rFactory),
		battleInteractor: NewBattleInteractor(dFactory, rFactory),
	}
}

func (f factory) LoginInteractor() LoginInteractor {
	return f.loginInteractor
}

func (f factory) BattleInteractor() BattleInteractor {
	return f.battleInteractor
}
