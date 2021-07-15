package infrastructures

import (
	"github.com/swallowarc/tictactoe_battle_backend/internal/infrastructures/env"
	"github.com/swallowarc/tictactoe_battle_backend/internal/infrastructures/redis"
	"github.com/swallowarc/tictactoe_battle_backend/internal/interface_adapters/gateways"
)

type (
	factory struct {
		redisClient gateways.MemDBClient
	}
)

func NewFactory() gateways.Factory {
	return &factory{
		redisClient: redis.NewRedisClient(env.Redis),
	}
}

func (f factory) MemDBClient() gateways.MemDBClient {
	return f.redisClient
}
