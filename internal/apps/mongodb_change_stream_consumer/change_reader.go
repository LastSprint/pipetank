package mongodbchangestreamconsumer

import (
	"context"
	"log/slog"
	"sync"
	"time"
	"weak"

	"github.com/LastSprint/pipetank/internal/repo"
	"github.com/LastSprint/pipetank/pkg/common"
	"github.com/LastSprint/pipetank/pkg/mdb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type store interface {
	SaveChangeStreamToken(context.Context, string, bson.Raw) error
	WatchRawEvents(context.Context, common.CallbackFailable[repo.RawEventWatchModel]) error
}

type handler interface {
	HandleEvents(context.Context, []repo.Event) error
}

type ErrorHandlingStrategy int

const (
	LogAndSkip ErrorHandlingStrategy = 0
	SendToDLQ  ErrorHandlingStrategy = 1
	Fail       ErrorHandlingStrategy = 2

	TickFrequency time.Duration = time.Millisecond * 100
)

type Consumer struct {
	repo    store
	handler handler

	buffered      []repo.Event
	bufferMx      *sync.RWMutex
	tokenProvider weak.Pointer[mdb.ResumeTokenProvider]

	errorHandlingStrategy ErrorHandlingStrategy
	maxBufferSize         int
	bufferCleanUpPeriod   time.Duration
	key                   string
}

func New(
	r store,
	h handler,
	errorHandlingStrategy ErrorHandlingStrategy,
	maxBufferSize int,
	bufferCleanUpPeriod time.Duration,
	key string,
) *Consumer {
	return &Consumer{
		repo:                  r,
		handler:               h,
		errorHandlingStrategy: errorHandlingStrategy,
		bufferCleanUpPeriod:   bufferCleanUpPeriod,
		maxBufferSize:         maxBufferSize,
		key:                   key,

		bufferMx: &sync.RWMutex{},
		buffered: make([]repo.Event, 0, maxBufferSize),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	// ctx, cancelFn := context.WithCancelCause()
	go func() {
		for {
			err := c.bufferFlusher(ctx)
			if err == nil {
				continue
			}
			err = c.onError(ctx, err)
			if err == nil {
				continue
			}

			// TODO: handle error respecting the strategy
			// cancelFn(err)
		}
	}()

	return c.repo.WatchRawEvents(
		ctx,
		func(ctx context.Context, event repo.RawEventWatchModel) error {
			err := c.handleEventAction(ctx, event)
			if err == nil {
				return nil
			}

			err = c.onError(ctx, err)
			if err == nil {
				return nil
			}

			return err
		},
	)
}

func (c *Consumer) handleEventAction(ctx context.Context, event repo.RawEventWatchModel) error {
	c.tokenProvider = weak.Make[mdb.ResumeTokenProvider](&event.Token)
	err := c.handleRecord(ctx, event.Record)
	if err != nil {
		return c.onError(ctx, err)
	}

	return nil
}

func (c *Consumer) onError(ctx context.Context, err error) error {
	switch c.errorHandlingStrategy {
	case Fail:
		return err
	case LogAndSkip:
		slog.WarnContext(ctx, "failed to process with error", slog.Any("error", err))
	case SendToDLQ:
		panic("not implemented")
	}

	slog.Warn(
		"Error handling strategy is incorrect; Must be 0, 1 or 2",
		slog.Any("set_strategy", c.errorHandlingStrategy),
	)

	return nil
}

func (c *Consumer) onSuccess(ctx context.Context, events []repo.Event) error {
	eventsIDs := make([]string, 0, len(events))
	for _, event := range events {
		eventsIDs = append(eventsIDs, event.ID.String())
	}

	slog.InfoContext(ctx, "processing events", slog.Any("events", eventsIDs))

	tp := c.tokenProvider.Value()

	if tp == nil {
		slog.Error("Weak ref to token provider is nil, cannot save token")
		return nil
	}

	ttp := *tp

	return c.repo.SaveChangeStreamToken(ctx, c.key, ttp.ResumeToken())
}

func (c *Consumer) handleRecord(_ context.Context, event repo.Event) error { //nolint:unparam
	c.bufferMx.Lock()
	c.buffered = append(c.buffered, event)
	c.bufferMx.Unlock()

	return nil
}

func (c *Consumer) bufferFlusher(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(TickFrequency):
			c.bufferMx.RLock()
			count := len(c.buffered)
			c.bufferMx.RUnlock()

			if count < c.maxBufferSize {
				break
			}

			err := c.flushBuffer(ctx)
			if err != nil {
				return err
			}

		case <-time.After(c.bufferCleanUpPeriod):
			err := c.flushBuffer(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Consumer) flushBuffer(ctx context.Context) error {
	if len(c.buffered) == 0 {
		return nil
	}

	// POTENTIAL OPTIMISATION
	// Use rolling buffers, in order to prevent data loss between copying and handling data
	// when we start flushing - this object must start writing data to another buffer
	c.bufferMx.Lock()
	defer c.bufferMx.Unlock()
	cp := make([]repo.Event, len(c.buffered))
	copy(cp, c.buffered)

	err := c.handler.HandleEvents(ctx, cp)
	if err != nil {
		return err
	}

	c.buffered = c.buffered[:0]

	return c.onSuccess(ctx, cp)
}
