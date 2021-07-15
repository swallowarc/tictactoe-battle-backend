package management_state

type (
	State int
)

const (
	Meeting State = 1
	Error   State = iota + 1
	PlayerATurn
	PlayerAPicked
	PlayerBTurn
	PlayerBPicked
	PlayerAWin
	PlayerBWin
)
