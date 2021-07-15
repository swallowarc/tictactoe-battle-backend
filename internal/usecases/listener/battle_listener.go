package listener

import (
	"context"

	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/management_state"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/room"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/ports"
	"golang.org/x/xerrors"
)

type (
	battleListener struct {
		roomID        room.ID
		loginID       string
		battleRepo    ports.BattleRepository
		lastMessageID string
	}
)

var LeftError = xerrors.New("already left the room")

func NewBattleListener(roomID room.ID, loginID string, battleRepo ports.BattleRepository) ports.BattleListener {
	return &battleListener{
		roomID:     roomID,
		loginID:    loginID,
		battleRepo: battleRepo,
	}
}

func (l *battleListener) Listen(ctx context.Context) (*tictactoe_battle.BattleSituation, error) {
	var (
		newMsgID string
		bt       *battle.Battle
		err      error
	)

	isExists, err := l.battleRepo.IsExistsInRoom(ctx, l.roomID, l.loginID)
	if err != nil {
		return nil, xerrors.Errorf("failed to IsExistsInRoom: %w", err)
	}
	if !isExists {
		return nil, LeftError
	}

	if l.lastMessageID == "" {
		newMsgID, bt, err = l.battleRepo.ReadStreamLatest(ctx, l.roomID)
		if err != nil {
			return nil, xerrors.Errorf("failed to ReadStreamLatest: %w", err)
		}
	} else {
		newMsgID, bt, err = l.battleRepo.ReadStream(ctx, l.roomID, l.lastMessageID)
		if err != nil {
			return nil, xerrors.Errorf("failed to ReadStream: %w", err)
		}
	}

	l.lastMessageID = newMsgID

	ret, err := l.convert(bt)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (l *battleListener) convert(b *battle.Battle) (*tictactoe_battle.BattleSituation, error) {
	ret := &tictactoe_battle.BattleSituation{
		RoomId:         b.RoomID.String(),
		PlayerAId:      b.PlayerAID,
		PlayerBId:      b.PlayerBID,
		PickedPosition: b.PickedPosition,
		PickedPiece:    b.PickedPiece,
		Field:          b.Field,
		WinLine:        b.WinLine,
	}

	switch l.loginID {
	case b.PlayerBID:
		ret.Player = tictactoe_battle.Player_PLAYER_B
		ret.Holding = b.PlayerBHolding

		switch b.State {
		case management_state.PlayerATurn:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_OPPONENT_TURN
		case management_state.PlayerAPicked:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_OPPONENT_TURN_PICKED
		case management_state.PlayerAWin:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_LOSE
		case management_state.PlayerBTurn:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_PLAYER_TURN
		case management_state.PlayerBPicked:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_PLAYER_TURN_PICKED
		case management_state.PlayerBWin:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_WIN
		case management_state.Meeting:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_MEETING
		case management_state.Error:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_ERROR
		default:
			return nil, xerrors.Errorf("unexpected management state: %s", b.State)
		}
	default:
		if l.loginID == b.PlayerAID {
			ret.Player = tictactoe_battle.Player_PLAYER_A
			ret.Holding = b.PlayerAHolding
		} else {
			ret.Player = tictactoe_battle.Player_PLAYER_AUDIENCE
			ret.Holding = &tictactoe_battle.Holding{}
		}

		switch b.State {
		case management_state.PlayerATurn:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_PLAYER_TURN
		case management_state.PlayerAPicked:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_PLAYER_TURN_PICKED
		case management_state.PlayerAWin:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_WIN
		case management_state.PlayerBTurn:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_OPPONENT_TURN
		case management_state.PlayerBPicked:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_OPPONENT_TURN_PICKED
		case management_state.PlayerBWin:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_LOSE
		case management_state.Meeting:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_MEETING
		case management_state.Error:
			ret.State = tictactoe_battle.BattleState_BATTLE_STATE_ERROR
		default:
			return nil, xerrors.Errorf("unexpected management state: %s", b.State)
		}
	}

	return ret, nil
}
