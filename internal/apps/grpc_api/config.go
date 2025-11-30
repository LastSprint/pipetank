package grpc_api

import "github.com/caarlos0/env/v11"

type config struct {
	Port int `env:"PORT" default:"50051"`
}

func parseConfig() (config, error) {
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
