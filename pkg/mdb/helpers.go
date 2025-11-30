package mdb

import (
	"context"

	"github.com/LastSprint/pipetank/pkg/observability/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (c *Client) Close(ctx context.Context) error {
	err := c.mdbClient.Disconnect(ctx)
	if err != nil {
		return errors.NewTErr(ctx, err, errors.ErrInternal)
	}
	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	err := c.mdbClient.Ping(ctx, nil)
	if err != nil {
		return errors.NewTErr(ctx, err, errors.ErrInternal)
	}
	return nil
}

func (c *Client) DB() *mongo.Database {
	return c.mdbClient.Database(c.dbName)
}

func (c *Client) Client() *mongo.Client {
	return c.mdbClient
}
