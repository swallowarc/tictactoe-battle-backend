package battle

import (
	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/management_state"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/room"
)

type (
	Battle struct {
		RoomID         room.ID                        `json:"room_id"`
		PlayerAID      string                         `json:"player_aid"`
		PlayerAHolding *tictactoe_battle.Holding      `json:"player_a_holding"`
		PlayerBID      string                         `json:"player_bid"`
		PlayerBHolding *tictactoe_battle.Holding      `json:"player_b_holding"`
		State          management_state.State         `json:"state"`
		PickedPosition tictactoe_battle.Position      `json:"picked_position"`
		PickedPiece    tictactoe_battle.Piece         `json:"picked_piece"`
		Field          []*tictactoe_battle.PieceStack `json:"field"`
		WinLine        tictactoe_battle.WinLine       `json:"win_line"`
	}
)
