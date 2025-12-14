package grpc_api

import (
	"context"
	"fmt"
	"net"

	"github.com/LastSprint/pipetank/internal/repo"
	raweventsconsumer "github.com/LastSprint/pipetank/internal/reusable/raw_events_consumer"
	"github.com/LastSprint/pipetank/pkg/client/proto"
	"github.com/LastSprint/pipetank/pkg/mdb"
	"github.com/LastSprint/pipetank/pkg/utils"
	"google.golang.org/grpc"
)

func Run(ctx context.Context) error {
	cfg, err := parseConfig()
	if err != nil {
		return fmt.Errorf("failed to parse Config: %w", err)
	}

	return RunWithConfig(ctx, cfg)
}

func RunWithConfig(ctx context.Context, cfg Config) error {
	mdbClinet, err := mdb.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create mongodb client: %w", err)
	}

	rep, err := repo.NewRepo(ctx, mdbClinet, utils.UTCClock())
	if err != nil {
		return err
	}

	srv := raweventsconsumer.NewService(rep)

	hdnls := newHandlers(srv)
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", cfg.Port, err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterAPIServer(grpcServer, hdnls)

	return utils.DieWithGrace(
		ctx,
		func(ctx context.Context) error {
			return grpcServer.Serve(listener)
		},
		func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	)
}
