package room

import (
	"fmt"
	"time"

	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/random"
)

const (
	idPattern       = "1234567890"
	idKeyPrefix     = "tic_tac_toe_id"
	memberKeyPrefix = "tic_tac_toe_member"
	streamKeyPrefix = "tic_tac_toe_stream"
)

const (
	TimeoutDuration = 15 * time.Minute
)

type (
	ID string
)

func NewID() ID {
	return ID(random.RandString6ByParam(5, idPattern))
}

func (id ID) IDKey() string {
	return fmt.Sprintf("%s:%s", idKeyPrefix, id)
}

func (id ID) MemberKey() string {
	return fmt.Sprintf("%s:%s", memberKeyPrefix, id)
}

func (id ID) StreamKey() string {
	return fmt.Sprintf("%s:%s", streamKeyPrefix, id)
}

func (id ID) String() string {
	return string(id)
}
