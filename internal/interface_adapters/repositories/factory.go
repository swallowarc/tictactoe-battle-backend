package repositories

import (
	"github.com/swallowarc/tictactoe_battle_backend/internal/interface_adapters/gateways"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/ports"
)

type (
	factory struct {
		loginRepository  ports.LoginRepository
		battleRepository ports.BattleRepository
	}
)

func NewFactory(gwFactory gateways.Factory) ports.RepositoriesFactory {
	return &factory{
		loginRepository:  NewLoginRepository(gwFactory),
		battleRepository: NewBattleRepository(gwFactory),
	}
}

func (f *factory) LoginRepository() ports.LoginRepository {
	return f.loginRepository
}

func (f *factory) BattleRepository() ports.BattleRepository {
	return f.battleRepository
}
