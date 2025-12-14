package utils

import (
	"context"
	"errors"
	"os/signal"
	"syscall"
)

func DieWithGrace(
	ctx context.Context,
	activity func(ctx context.Context) error,
	onDone func(ctx context.Context) error,
) error {
	newCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	activityDoneChan := make(chan error, 1)

	go func() {
		activityDoneChan <- activity(newCtx)
	}()

	select {
	case <-newCtx.Done():
		return onDone(ctx)
	case err := <-activityDoneChan:
		return errors.Join(err, onDone(ctx))
	}
}
