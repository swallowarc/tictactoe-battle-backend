package interactors

import (
	"context"

	"github.com/swallowarc/tictactoe-battle-proto/pkg/tictactoe_battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/exceptions"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/room"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/listener"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/ports"
	"golang.org/x/xerrors"
)

type (
	battleInteractor struct {
		battleRule battle.Rule
		battleRepo ports.BattleRepository
	}
)

func NewBattleInteractor(dFactory domains.Factory, rFactory ports.RepositoriesFactory) BattleInteractor {
	return &battleInteractor{
		battleRule: dFactory.BattleRule(),
		battleRepo: rFactory.BattleRepository(),
	}
}

func (bi *battleInteractor) Create(ctx context.Context) (room.ID, error) {
	roomID, err := bi.battleRepo.Create(ctx, bi.battleRule.OpenBattle())
	if err != nil {
		return "", xerrors.Errorf("failed to Create: %w", err)
	}

	return roomID, nil
}

func (bi *battleInteractor) CanEnter(ctx context.Context, roomID room.ID, _ string) (bool, error) {
	_, _, err := bi.battleRepo.ReadStreamLatest(ctx, roomID)
	if err != nil {
		if exceptions.IsNotFoundError(err) {
			return false, nil
		}
		return false, xerrors.Errorf("failed to ReadStreamLatest: %w", err)
	}

	return true, nil
}

func (bi *battleInteractor) Enter(ctx context.Context, roomID room.ID, loginID string) (ports.BattleListener, error) {
	if err := bi.battleRepo.Enter(ctx, roomID, loginID); err != nil {
		return nil, xerrors.Errorf("failed to Enter: %w", err)
	}

	_, _, err := bi.battleRepo.ReadStreamLatest(ctx, roomID)
	if err != nil {
		return nil, xerrors.Errorf("failed to ReadStreamLatest: %w", err)
	}

	return listener.NewBattleListener(roomID, loginID, bi.battleRepo), nil
}

func (bi *battleInteractor) Declaration(ctx context.Context, roomID room.ID, loginID string) error {
	_, bt, err := bi.battleRepo.ReadStreamLatest(ctx, roomID)
	if err != nil {
		return xerrors.Errorf("failed to ReadStreamLatest: %w", err)
	}

	if err := bi.battleRule.Declaration(bt, loginID); err != nil {
		return xerrors.Errorf("failed to Declaration: %w", err)
	}

	if err := bi.battleRepo.Update(ctx, bt); err != nil {
		return xerrors.Errorf("failed to Update: %w", err)
	}

	return nil
}

func (bi *battleInteractor) Leave(ctx context.Context, roomID room.ID, loginID string) error {
	if err := bi.battleRepo.Leave(ctx, roomID, loginID); err != nil {
		return xerrors.Errorf("failed to Leave: %w", err)
	}

	// 最終退出者だった場合はroomを削除する
	members, err := bi.battleRepo.ListMembers(ctx, roomID)
	if err != nil {
		return xerrors.Errorf("failed to ListMembers: %w", err)
	}
	if len(members) != 0 {
		return nil
	}
	if err := bi.battleRepo.Delete(ctx, roomID); err != nil {
		return xerrors.Errorf("failed to Delete: %w", err)
	}
	return nil
}

func (bi *battleInteractor) Attack(ctx context.Context, roomID room.ID, player tictactoe_battle.Player, position tictactoe_battle.Position, pieceSize tictactoe_battle.Piece) error {
	_, b, err := bi.battleRepo.ReadStreamLatest(ctx, roomID)
	if err != nil {
		return xerrors.Errorf("failed to ReadStreamLatest: %w", err)
	}

	if err := bi.battleRule.Attack(b, player, position, pieceSize); err != nil {
		return xerrors.Errorf("failed to bt Attack: %w", err)
	}

	if err := bi.battleRepo.Update(ctx, b); err != nil {
		return xerrors.Errorf("failed to bt Update: %w", err)
	}

	return nil
}

func (bi *battleInteractor) Pick(ctx context.Context, roomID room.ID, player tictactoe_battle.Player, position tictactoe_battle.Position, pieceSize tictactoe_battle.Piece) error {
	_, b, err := bi.battleRepo.ReadStreamLatest(ctx, roomID)
	if err != nil {
		return xerrors.Errorf("failed to ReadStreamLatest: %w", err)
	}

	if err := bi.battleRule.Pick(b, player, position, pieceSize); err != nil {
		return xerrors.Errorf("failed to bt Attack: %w", err)
	}

	if err := bi.battleRepo.Update(ctx, b); err != nil {
		return xerrors.Errorf("failed to bt Update: %w", err)
	}

	return nil
}

func (bi *battleInteractor) Reset(ctx context.Context, roomID room.ID) error {
	_, b, err := bi.battleRepo.ReadStreamLatest(ctx, roomID)
	if err != nil {
		return xerrors.Errorf("failed to ReadStreamLatest: %w", err)
	}
	bi.battleRule.Reset(b)

	if err := bi.battleRepo.Update(ctx, b); err != nil {
		return xerrors.Errorf("failed to battle Update: %w", err)
	}

	return nil
}
