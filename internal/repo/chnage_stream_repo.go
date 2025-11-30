package repo

import (
	"context"

	"github.com/LastSprint/pipetank/pkg/observability/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

const (
	NamespaceChangeStreamTokenStorage = "change_stream_token_storage"
)

func (r *Repo) SaveChangeStreamToken(ctx context.Context, key string, token bson.Raw) error {
	colOpts := options.Collection().
		SetWriteConcern(writeconcern.W1()).
		SetReadPreference(readpref.SecondaryPreferred())

	opOpts := options.FindOneAndUpdate().SetUpsert(true)

	opResult := r.client.DB().
		Collection(NamespaceChangeStreamTokenStorage, colOpts).
		FindOneAndUpdate(ctx, bson.M{"key": key}, bson.M{"token": bson.M{"$set": token}}, opOpts)

	err := opResult.Err()
	if err != nil {
		return errors.NewTErr(ctx, err, errors.ErrInternal)
	}

	return nil
}
