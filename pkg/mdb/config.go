package mdb

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DSN               string        `env:"MDB_DSN,required"`
	DB                string        `env:"MDB_DB,required"`
	MaxConnections    int           `env:"MDB_MAX_CONNECTIONS,required"`
	AppName           string        `env:"APP_NAME,required"`
	ConnectionTimeout time.Duration `env:"MDB_CONNECTION_TIMEOUT"       envDefault:"10s"`
	MaxConnecting     uint64        `env:"MDB_MAX_CONNECTING"           envDefault:"2"`
	MaxPoolSize       uint64        `env:"MDB_MAX_POOL_SIZE"            envDefault:"10"`
	MinPoolSize       uint64        `env:"MDB_MIN_POOL_SIZE"            envDefault:"1"`
}

func ParseConfig() (Config, error) {
	return env.ParseAs[Config]()
}
