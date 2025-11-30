package mongodbchangestreamconsumer

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type config struct {
	HealthCheckAddr string `env:"HEALTH_CHECK_ADDR"`

	ErrorHandlingStrategy ErrorHandlingStrategy `env:"ERROR_HANDLING_STRATEGY,default:0"`
	MaxBufferSize         int                   `env:"MAX_BUFFER_SIZE,default:1000"`
	BufferCleanUpPeriod   time.Duration         `env:"BUFFER_CLEANUP_PERIOD,default:5s"`
	ConsumerKey           string                `env:"CONSUMER_KEY,default:all_in_one"`
}

func parseConfig() (config, error) {
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
