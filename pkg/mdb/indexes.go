package mdb

import (
	"context"
	"errors"
	"log/slog"

	octx "github.com/LastSprint/pipetank/pkg/observability/ctx"
	oerrs "github.com/LastSprint/pipetank/pkg/observability/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MGOIndex struct {
	Name               string `bson:"name"`
	ExpireAfterSeconds *int32 `bson:"expireAfterSeconds,omitempty"`
	Version            *int32 `bson:"v,omitempty"`
}

// GetIndexByName tries to find an index with a given name in a given collection
// Errors:
// - oerrs.ErrInternal: on any internal error
// - oerrs.ErrNotFound: if the index is not found.
func (c *Client) GetIndexByName(
	ctx context.Context,
	collection, indexName string,
) (index MGOIndex, err error) {
	indexesCur, err := c.DB().Collection(collection).Indexes().List(ctx)
	if err != nil {
		return index, oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
	}

	var indexes []MGOIndex

	err = indexesCur.All(ctx, &indexes)
	if err != nil {
		return index, oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
	}

	for _, idx := range indexes {
		if idx.Name == indexName {
			return idx, nil
		}
	}

	return index, oerrs.NewTErrf(
		ctx,
		"index %s not found in col %s: %w",
		indexName,
		collection,
		oerrs.ErrNotFound,
	)
}

// CreateIndexes creates indexes in a given collection
// Errors:
// - oerrs.ErrInternal: on any error.
func (c *Client) CreateIndexes(
	ctx context.Context,
	collectionName string,
	indexes []mongo.IndexModel,
) error {
	_, err := c.DB().Collection(collectionName).
		Indexes().
		CreateMany(
			ctx,
			indexes,
		)
	if err != nil {
		return oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
	}

	return nil
}

// CreateOrUpdateTTLIndex creates TTL index if it does not exist or updates it if it exists (via dropping)
// Errors:
// - oerrs.ErrInternal: on any error.
func (c *Client) CreateOrUpdateTTLIndex(
	ctx context.Context,
	collection, indexName string,
	expireAfterSeconds int32,
	newIndex mongo.IndexModel,
) error {
	index, err := c.GetIndexByName(ctx, collection, indexName)

	if errors.Is(err, oerrs.ErrNotFound) {
		return c.CreateIndexes(
			ctx,
			collection,
			[]mongo.IndexModel{newIndex},
		)
	}

	if err != nil {
		return err
	}

	if index.ExpireAfterSeconds == nil {
		return oerrs.NewTErrf(ctx, "index %s is not TTL index: %w", indexName, oerrs.ErrInternal)
	}

	if *index.ExpireAfterSeconds == expireAfterSeconds {
		octx.Logger(ctx).
			WithGroup("CreateOrUpdateTTLIndex").
			With(slog.Any("index", index)).
			Debug("skipped index update")

		return nil
	}

	// here we need to update an index

	err = c.DB().Collection(collection).Indexes().DropOne(ctx, index.Name)
	if err != nil {
		return oerrs.NewTErrf(ctx, "failed to drop index %s: %w", index.Name, oerrs.ErrInternal)
	}

	return c.CreateIndexes(
		ctx,
		collection,
		[]mongo.IndexModel{newIndex},
	)
}
