//go:generate mockgen -source=$GOFILE -destination=../../tests/mocks/domains/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package battle

import (
	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/exceptions"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/management_state"
	"golang.org/x/xerrors"
)

type (
	Rule interface {
		OpenBattle() *Battle
		Declaration(b *Battle, playerID string) error
		Attack(b *Battle, player tictactoe_battle.Player, pos tictactoe_battle.Position, size tictactoe_battle.Piece) error
		Pick(b *Battle, player tictactoe_battle.Player, pos tictactoe_battle.Position, size tictactoe_battle.Piece) error
		Reset(b *Battle)
	}

	rule struct{}
)

func NewRule() Rule {
	return &rule{}
}

func (r *rule) OpenBattle() *Battle {
	return &Battle{
		RoomID:         "",
		PlayerAID:      "",
		PlayerAHolding: newDefaultHolding(),
		PlayerBID:      "",
		PlayerBHolding: newDefaultHolding(),
		State:          management_state.Meeting,
		PickedPosition: tictactoe_battle.Position_POSITION_UNDEFINED,
		PickedPiece:    tictactoe_battle.Piece_PIECE_UNKNOWN,
		Field:          newDefaultFiled(),
		WinLine:        tictactoe_battle.WinLine_WIN_LINE_UNKNOWN,
	}
}

func (r *rule) Declaration(b *Battle, playerID string) error {
	// 切断後の再接続を考慮
	if b.PlayerAID == playerID || b.PlayerBID == playerID {
		return nil
	}

	if b.PlayerAID == "" {
		b.PlayerAID = playerID
		return nil
	}

	if b.PlayerBID == "" {
		b.PlayerBID = playerID
		b.State = management_state.PlayerATurn
		return nil
	}

	return exceptions.NewPreConditionError("battle has already started")
}

func (r *rule) Attack(b *Battle, player tictactoe_battle.Player, pos tictactoe_battle.Position, size tictactoe_battle.Piece) error {
	// check turn
	var valid bool
	switch b.State {
	case management_state.PlayerATurn, management_state.PlayerAPicked:
		valid = player == tictactoe_battle.Player_PLAYER_A
	case management_state.PlayerBTurn, management_state.PlayerBPicked:
		valid = player == tictactoe_battle.Player_PLAYER_B
	default:
		// invalid
	}
	if !valid {
		return xerrors.Errorf("it is not %s's turn: %s", player, b.State)
	}

	// check dest position
	if pos == tictactoe_battle.Position_POSITION_UNDEFINED {
		return xerrors.Errorf("selected unexpected position: %s", pos)
	}

	// check placement in picked field
	if b.PickedPosition == pos {
		return xerrors.New("cannot be relocated to the field from which it was picked")
	}

	// check placement from picked size
	if b.PickedPiece != tictactoe_battle.Piece_PIECE_UNKNOWN && b.PickedPiece != size {
		return xerrors.New("only the selected piece can be rearranged")
	}

	// check stack & put piece
	const (
		MissingPiecesMsg = "missing pieces on holding: %s"
		sizeInvalidMsg   = "pieces larger than the specified size have been placed. pos: %s, size: %s"
		stackInvalidMsg  = "cannot stack it on own piece"
	)

	var holding *tictactoe_battle.Holding
	switch player {
	case tictactoe_battle.Player_PLAYER_A:
		holding = b.PlayerAHolding
	default:
		holding = b.PlayerBHolding
	}

	stack := b.Field[pos]
	switch size {
	case tictactoe_battle.Piece_PIECE_S:
		if holding.S <= 0 {
			return xerrors.Errorf(MissingPiecesMsg, size)
		}
		if stack.L != tictactoe_battle.Player_PLAYER_UNKNOWN ||
			stack.M != tictactoe_battle.Player_PLAYER_UNKNOWN ||
			stack.S != tictactoe_battle.Player_PLAYER_UNKNOWN {
			return xerrors.Errorf(sizeInvalidMsg, pos, size)
		}

		stack.S = player
		holding.S--

	case tictactoe_battle.Piece_PIECE_M:
		if holding.M <= 0 {
			return xerrors.Errorf(MissingPiecesMsg, size)
		}
		if stack.L != tictactoe_battle.Player_PLAYER_UNKNOWN ||
			stack.M != tictactoe_battle.Player_PLAYER_UNKNOWN {
			return xerrors.Errorf(sizeInvalidMsg, pos, size)
		}
		if stack.S == player {
			return xerrors.New(stackInvalidMsg)
		}

		stack.M = player
		holding.M--

	case tictactoe_battle.Piece_PIECE_L:
		if holding.L <= 0 {
			return xerrors.Errorf(MissingPiecesMsg, size)
		}
		if stack.L != tictactoe_battle.Player_PLAYER_UNKNOWN {
			return xerrors.Errorf(sizeInvalidMsg, pos, size)
		}
		if stack.M == player ||
			stack.M == tictactoe_battle.Player_PLAYER_UNKNOWN && stack.S == player {
			return xerrors.New(stackInvalidMsg)
		}

		stack.L = player
		holding.L--
	}

	// reset picked state
	b.PickedPosition = tictactoe_battle.Position_POSITION_UNDEFINED
	b.PickedPiece = tictactoe_battle.Piece_PIECE_UNKNOWN

	r.judgment(b)

	// turn change
	switch b.State {
	case management_state.PlayerAWin, management_state.PlayerBWin:
		return nil
	case management_state.PlayerATurn, management_state.PlayerAPicked:
		b.State = management_state.PlayerBTurn
	case management_state.PlayerBTurn, management_state.PlayerBPicked:
		b.State = management_state.PlayerATurn
	}

	return nil
}

func (r *rule) Pick(b *Battle, player tictactoe_battle.Player, pos tictactoe_battle.Position, size tictactoe_battle.Piece) error {
	var valid bool
	switch b.State {
	case management_state.PlayerATurn:
		valid = player == tictactoe_battle.Player_PLAYER_A
		b.State = management_state.PlayerAPicked
	case management_state.PlayerBTurn:
		valid = player == tictactoe_battle.Player_PLAYER_B
		b.State = management_state.PlayerBPicked
	default:
		// invalid
	}
	if !valid {
		return xerrors.Errorf("it is not %s's turn: %s", player, b.State)
	}

	// check stack & pick
	const (
		largerPiecesMsg  = "larger pieces are placed. pos: %s"
		playerInvalidMsg = "it's not the player's piece. player: %s"
	)
	stack := b.Field[pos]
	var holding *tictactoe_battle.Holding
	switch player {
	case tictactoe_battle.Player_PLAYER_A:
		holding = b.PlayerAHolding
	default:
		holding = b.PlayerBHolding
	}

	switch size {
	case tictactoe_battle.Piece_PIECE_S:
		if stack.L != tictactoe_battle.Player_PLAYER_UNKNOWN ||
			stack.M != tictactoe_battle.Player_PLAYER_UNKNOWN {
			return xerrors.Errorf(largerPiecesMsg, pos)
		}
		if player != stack.S {
			return xerrors.Errorf(playerInvalidMsg, stack.S)
		}

		stack.S = tictactoe_battle.Player_PLAYER_UNKNOWN
		holding.S++

	case tictactoe_battle.Piece_PIECE_M:
		if stack.L != tictactoe_battle.Player_PLAYER_UNKNOWN {
			return xerrors.Errorf(largerPiecesMsg, pos)
		}
		if player != stack.M {
			return xerrors.Errorf(playerInvalidMsg, stack.M)
		}

		stack.M = tictactoe_battle.Player_PLAYER_UNKNOWN
		holding.M++

	case tictactoe_battle.Piece_PIECE_L:
		if player != stack.L {
			return xerrors.Errorf(playerInvalidMsg, stack.L)
		}

		stack.L = tictactoe_battle.Player_PLAYER_UNKNOWN
		holding.L++
	}

	b.PickedPosition = pos
	b.PickedPiece = size

	r.judgment(b)

	return nil
}

func (r *rule) Reset(b *Battle) {
	b.State = management_state.Meeting
	b.PlayerAID = ""
	b.PlayerBID = ""
	b.PlayerAHolding = newDefaultHolding()
	b.PlayerBHolding = newDefaultHolding()
	b.Field = newDefaultFiled()
	b.PickedPosition = tictactoe_battle.Position_POSITION_UNDEFINED
	b.PickedPiece = tictactoe_battle.Piece_PIECE_UNKNOWN
	b.WinLine = tictactoe_battle.WinLine_WIN_LINE_UNKNOWN
}

func (r *rule) judgment(b *Battle) {
	field := [][]tictactoe_battle.Player{
		{stackOwner(b.Field[0]), stackOwner(b.Field[3]), stackOwner(b.Field[6])},
		{stackOwner(b.Field[1]), stackOwner(b.Field[4]), stackOwner(b.Field[7])},
		{stackOwner(b.Field[2]), stackOwner(b.Field[5]), stackOwner(b.Field[8])},
	}

	winPlayer := func(player tictactoe_battle.Player) management_state.State {
		if player == tictactoe_battle.Player_PLAYER_A {
			return management_state.PlayerAWin
		}
		return management_state.PlayerBWin
	}

	if target := field[0][0]; target != tictactoe_battle.Player_PLAYER_UNKNOWN {
		if target == field[1][0] && target == field[2][0] {
			b.WinLine = tictactoe_battle.WinLine_WIN_LINE_1
			b.State = winPlayer(target)
			return
		}
		if target == field[0][1] && target == field[0][2] {
			b.WinLine = tictactoe_battle.WinLine_WIN_LINE_4
			b.State = winPlayer(target)
			return
		}
		if target == field[1][1] && target == field[2][2] {
			b.WinLine = tictactoe_battle.WinLine_WIN_LINE_7
			b.State = winPlayer(target)
			return
		}
	}

	if target := field[0][1]; target != tictactoe_battle.Player_PLAYER_UNKNOWN {
		if target == field[1][1] && target == field[2][1] {
			b.WinLine = tictactoe_battle.WinLine_WIN_LINE_2
			b.State = winPlayer(target)
			return
		}
	}

	if target := field[0][2]; target != tictactoe_battle.Player_PLAYER_UNKNOWN {
		if target == field[1][2] && target == field[2][2] {
			b.WinLine = tictactoe_battle.WinLine_WIN_LINE_3
			b.State = winPlayer(target)
			return
		}
	}

	if target := field[1][0]; target != tictactoe_battle.Player_PLAYER_UNKNOWN {
		if target == field[1][1] && target == field[1][2] {
			b.WinLine = tictactoe_battle.WinLine_WIN_LINE_5
			b.State = winPlayer(target)
			return
		}
	}

	if target := field[2][0]; target != tictactoe_battle.Player_PLAYER_UNKNOWN {
		if target == field[2][1] && target == field[2][2] {
			b.WinLine = tictactoe_battle.WinLine_WIN_LINE_6
			b.State = winPlayer(target)
			return
		}
		if target == field[1][1] && target == field[0][2] {
			b.WinLine = tictactoe_battle.WinLine_WIN_LINE_8
			b.State = winPlayer(target)
			return
		}
	}
}

func stackOwner(s *tictactoe_battle.PieceStack) tictactoe_battle.Player {
	if s.L != tictactoe_battle.Player_PLAYER_UNKNOWN {
		return s.L
	}
	if s.M != tictactoe_battle.Player_PLAYER_UNKNOWN {
		return s.M
	}
	if s.S != tictactoe_battle.Player_PLAYER_UNKNOWN {
		return s.S
	}

	return tictactoe_battle.Player_PLAYER_UNKNOWN
}

func newDefaultHolding() *tictactoe_battle.Holding {
	return &tictactoe_battle.Holding{
		S: 2,
		M: 2,
		L: 2,
	}
}

func newDefaultFiled() []*tictactoe_battle.PieceStack {
	ret := make([]*tictactoe_battle.PieceStack, 9)
	for i := 0; i < 9; i++ {
		ret[i] = &tictactoe_battle.PieceStack{
			S: tictactoe_battle.Player_PLAYER_UNKNOWN,
			M: tictactoe_battle.Player_PLAYER_UNKNOWN,
			L: tictactoe_battle.Player_PLAYER_UNKNOWN,
		}
	}

	return ret
}
