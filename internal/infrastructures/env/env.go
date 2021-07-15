package env

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/mode"
	"github.com/swallowarc/tictactoe_battle_backend/internal/infrastructures/redis"
)

var (
	Server Config
	Redis  redis.Config
)

type (
	Config struct {
		RunMode mode.Mode `envconfig:"run_mode" default:"debug"`
		PORT    string    `envconfig:"grpc_port" default:"50051"`
	}
)

func init() {
	setup()
}

func setup() {
	check(envconfig.Process("", &Server))
	check(envconfig.Process("redis", &Redis))
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
