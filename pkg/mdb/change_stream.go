package mdb

import (
	"context"
	"log/slog"

	oerrs "github.com/LastSprint/pipetank/pkg/observability/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ResumeTokenProvider interface {
	ResumeToken() bson.Raw
}

func RunChangeStream[T any](
	ctx context.Context,
	c *Client,
	colName string,
	pipeline any,
	action func(ctx context.Context, token ResumeTokenProvider, doc T) error,
) error {
	return c.changeStream(
		ctx,
		colName,
		pipeline,
		func(ctx context.Context, token ResumeTokenProvider, doc bson.Raw) error {
			c, ok := doc.Lookup("fullDocument").DocumentOK()
			if !ok {
				return nil
			}

			var dt T

			err := bson.Unmarshal(c, &dt)
			if err != nil {
				return err
			}

			return action(ctx, token, dt)
		},
	)
}

func (c *Client) changeStream(
	ctx context.Context,
	colName string,
	pipeline any,
	action func(ctx context.Context, token ResumeTokenProvider, doc bson.Raw) error,
) error {
	cs, err := c.DB().
		Collection(colName).
		Watch(ctx, pipeline, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		return oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
	}

	defer func() {
		err = cs.Close(ctx)
		slog.Error("Failed to close change stream", "error", err)
	}()

	for cs.Next(ctx) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = cs.Err()
		if err != nil {
			return oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
		}

		err = action(ctx, cs, cs.Current)
		if err != nil {
			return err
		}
	}

	return nil
}
