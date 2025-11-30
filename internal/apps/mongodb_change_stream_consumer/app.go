package mongodbchangestreamconsumer

import (
	"context"
	"fmt"

	"github.com/LastSprint/pipetank/internal/repo"
	raweventsconsumer "github.com/LastSprint/pipetank/internal/reusable/raw_events_consumer"
	"github.com/LastSprint/pipetank/pkg/mdb"
	"github.com/LastSprint/pipetank/pkg/utils"
)

func Run(ctx context.Context) error {
	cfg, err := parseConfig()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	mdbClinet, err := mdb.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create mongodb client: %w", err)
	}

	rep, err := repo.NewRepo(ctx, mdbClinet, utils.UTCClock())
	if err != nil {
		return err
	}

	srv := raweventsconsumer.NewService(rep)

	consumer := New(
		rep,
		srv,
		cfg.ErrorHandlingStrategy,
		cfg.MaxBufferSize,
		cfg.BufferCleanUpPeriod,
		cfg.ConsumerKey,
	)

	return utils.DieWithGrace(
		ctx,
		func(ctx context.Context) error {
			return consumer.Start(ctx)
		},
		func(ctx context.Context) error {
			return mdbClinet.Close(ctx)
		},
	)
}
