package grpc_api

import "github.com/caarlos0/env/v11"

type Config struct {
	Port int `env:"PORT" default:"50051"`
}

func parseConfig() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
