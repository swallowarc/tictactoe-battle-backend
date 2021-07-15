package redis_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/exceptions"
	"github.com/swallowarc/tictactoe_battle_backend/internal/infrastructures/env"
	. "github.com/swallowarc/tictactoe_battle_backend/internal/infrastructures/redis"
)

func TestRedisClient(t *testing.T) {
	cli := NewRedisClient(env.Redis)
	ctx := context.Background()
	if err := cli.Ping(ctx); err != nil {
		t.Fatalf("failed to Ping: %v", err)
	}

	t.Run("SetNX, Get", func(t *testing.T) {
		key := uuid.NewString()
		val := uuid.NewString()

		const cacheLife = 10 * time.Second

		if err := cli.SetNX(ctx, key, val, cacheLife); err != nil {
			t.Fatalf("failed to SetNX: %v", err)
		}

		getVal, err := cli.Get(ctx, key)
		if err != nil {
			t.Fatalf("failed to Get: %v", err)
		}
		if val != getVal {
			t.Fatalf("wanted %s but got %s", val, getVal)
		}

		time.Sleep(cacheLife)

		getVal, err = cli.Get(ctx, key)
		if !exceptions.IsNotFoundError(err) {
			t.Fatalf("failed to Get: %v", err)
		}
		if getVal != "" {
			t.Fatalf("wanted empty but got %s", getVal)
		}
	})

	t.Run("SAdd, SRem, SMembers", func(t *testing.T) {
		key, val1, val2, val3 := uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString()

		if err := cli.SAdd(ctx, key, val1, val2, val3); err != nil {
			t.Fatalf("failed to SAdd: %v", err)
		}

		if err := cli.SRem(ctx, key, val1); err != nil {
			t.Fatalf("failed to SRem: %v", err)
		}

		members, err := cli.SMembers(ctx, key)
		if err != nil {
			t.Fatalf("failed to SMembers: %v", err)
		}
		sort.Slice(members, func(i, j int) bool {
			return members[i] < members[j]
		})
		want := []string{val2, val3}
		sort.Slice(want, func(i, j int) bool {
			return want[i] < want[j]
		})

		if diff := cmp.Diff(want, members); diff != "" {
			t.Fatalf(diff)
		}

		if err := cli.Del(ctx, key); err != nil {
			t.Fatalf("failed to Del: %v", err)
		}

		if r, err := cli.SMembers(ctx, key); err != nil {
			t.Fatalf("failed to SMembers: %v", err)
		} else if len(r) != 0 {
			t.Fatalf("got members %v, but want %v", r, []string{})
		}
	})
}
