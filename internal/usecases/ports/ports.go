package ports

type (
	RepositoriesFactory interface {
		LoginRepository() LoginRepository
		BattleRepository() BattleRepository
	}
)
