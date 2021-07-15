package repositories

import (
	"context"
	"encoding/json"

	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/exceptions"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/battle"
	"github.com/swallowarc/tictactoe_battle_backend/internal/domains/room"
	"github.com/swallowarc/tictactoe_battle_backend/internal/interface_adapters/gateways"
	"github.com/swallowarc/tictactoe_battle_backend/internal/usecases/ports"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

const (
	battleMessageKey = "tic_tac_toe_battle_message_key"
)

type (
	battleRepository struct {
		memDBCli gateways.MemDBClient
	}
)

func NewBattleRepository(gwFactory gateways.Factory) ports.BattleRepository {
	return &battleRepository{
		memDBCli: gwFactory.MemDBClient(),
	}
}

func (r *battleRepository) Create(ctx context.Context, battle *battle.Battle) (room.ID, error) {
	var roomID room.ID
	for {
		roomID = room.NewID()
		_, err := r.memDBCli.Get(ctx, roomID.IDKey())
		if err == nil {
			continue
		}
		if exceptions.IsNotFoundError(err) {
			break
		}
		return "", xerrors.Errorf("failed to memDBCli.Get: %w", err)
	}

	if err := r.memDBCli.SetNX(ctx, roomID.IDKey(), "", room.TimeoutDuration); err != nil {
		return "", xerrors.Errorf("failed to SetNX: %w", err)
	}

	battle.RoomID = roomID
	if err := r.Update(ctx, battle); err != nil { // UpdateでもStreamがなければ新規作成される
		return "", err
	}

	return roomID, nil
}

func (r *battleRepository) refreshRoomDuration(ctx context.Context, roomID room.ID) error {
	if _, err := r.memDBCli.Get(ctx, roomID.IDKey()); err != nil {
		return xerrors.Errorf("failed to Get room_id from memdb: %w", err)
	}

	eg := errgroup.Group{}

	eg.Go(func() error {
		if err := r.memDBCli.Expire(ctx, roomID.IDKey(), room.TimeoutDuration); err != nil {
			return xerrors.Errorf("failed to Expire room: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := r.memDBCli.Expire(ctx, roomID.MemberKey(), room.TimeoutDuration); err != nil {
			return xerrors.Errorf("failed to Expire member: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := r.memDBCli.Expire(ctx, roomID.StreamKey(), room.TimeoutDuration); err != nil {
			return xerrors.Errorf("failed to Expire stream: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (r *battleRepository) Update(ctx context.Context, battle *battle.Battle) error {
	if err := r.refreshRoomDuration(ctx, battle.RoomID); err != nil {
		return xerrors.Errorf("failed to refreshRoomDuration: %w", err)
	}

	jm, err := json.Marshal(battle)
	if err != nil {
		return xerrors.Errorf("failed to json.Marshal: %w", err)
	}

	if err := r.memDBCli.PublishStream(ctx, battle.RoomID.StreamKey(), map[string]interface{}{
		battleMessageKey: jm,
	}); err != nil {
		return xerrors.Errorf("failed to PublishStream: %w", err)
	}

	return nil
}

func (r *battleRepository) Enter(ctx context.Context, roomID room.ID, loginID string) error {
	if err := r.refreshRoomDuration(ctx, roomID); err != nil {
		return xerrors.Errorf("failed to refreshRoomDuration: %w", err)
	}

	if err := r.memDBCli.SAdd(ctx, roomID.MemberKey(), loginID); err != nil {
		return xerrors.Errorf("failed to SAdd room member from memdb: %w", err)
	}

	return nil
}

func (r *battleRepository) Leave(ctx context.Context, roomID room.ID, loginID string) error {
	if err := r.refreshRoomDuration(ctx, roomID); err != nil {
		return xerrors.Errorf("failed to refreshRoomDuration: %w", err)
	}

	if err := r.memDBCli.SRem(ctx, roomID.MemberKey(), loginID); err != nil {
		return xerrors.Errorf("failed to SRem from memdb: %w", err)
	}

	return nil
}

func (r *battleRepository) ReadStreamLatest(ctx context.Context, roomID room.ID) (string, *battle.Battle, error) {
	msgID, msg, err := r.memDBCli.ReadStreamLatest(ctx, roomID.StreamKey(), battleMessageKey)
	if err != nil {
		return "", nil, xerrors.Errorf("failed to ReadStreamLatest: %w", err)
	}

	result, err := unmarshal(msg)
	return msgID, result, err
}

func (r *battleRepository) ReadStream(ctx context.Context, roomID room.ID, messageID string) (string, *battle.Battle, error) {
	msgID, msg, err := r.memDBCli.ReadStream(ctx, roomID.StreamKey(), battleMessageKey, messageID)
	if err != nil {
		return "", nil, xerrors.Errorf("failed to ReadStream: %w", err)
	}

	result, err := unmarshal(msg)
	return msgID, result, err
}

func unmarshal(message string) (*battle.Battle, error) {
	var result battle.Battle
	if err := json.Unmarshal([]byte(message), &result); err != nil {
		return nil, xerrors.Errorf("failed to json unmarshal. err: %w, msg: %s", err, message)
	}
	return &result, nil
}

func (r *battleRepository) ListMembers(ctx context.Context, roomID room.ID) ([]string, error) {
	members, err := r.memDBCli.SMembers(ctx, roomID.MemberKey())
	if err != nil {
		return nil, xerrors.Errorf("failed to SMembers from memdb: %w", err)
	}

	return members, nil
}

func (r *battleRepository) IsExistsInRoom(ctx context.Context, roomID room.ID, loginID string) (bool, error) {
	memberIDs, err := r.ListMembers(ctx, roomID)
	if err != nil {
		return false, xerrors.Errorf("failed to ListMembers: %w", err)
	}

	for _, id := range memberIDs {
		if loginID == id {
			return true, nil
		}
	}

	return false, nil
}

func (r *battleRepository) Delete(ctx context.Context, roomID room.ID) error {
	eg := errgroup.Group{}
	eg.Go(func() error {
		return r.memDBCli.Del(ctx, roomID.IDKey())
	})
	eg.Go(func() error {
		return r.memDBCli.Del(ctx, roomID.MemberKey())
	})
	eg.Go(func() error {
		return r.memDBCli.Del(ctx, roomID.StreamKey())
	})
	if err := eg.Wait(); err != nil {
		return xerrors.Errorf("failed to Del from memdb: %w", err)
	}

	return nil
}
