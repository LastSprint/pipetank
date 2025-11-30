package repo

import (
	"context"
	"os"
	"strconv"

	oerrs "github.com/LastSprint/pipetank/pkg/observability/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	collectionNameRawEvents = "executions"

	idxNameRawEventsTTL = "ttl_raw_events"
)

func (r *Repo) createRawEventIndexes(ctx context.Context) error {
	ttlSec, err := strconv.ParseInt(os.Getenv("RAW_EVENTS_TTL_SECONDS"), 10, 32)
	if err != nil {
		return oerrs.NewTErr(ctx, err, oerrs.ErrBadInput)
	}

	err = r.client.CreateOrUpdateTTLIndex(
		ctx,
		collectionNameRawEvents,
		idxNameRawEventsTTL,
		int32(ttlSec),
		mongo.IndexModel{
			Keys: bson.M{RawEventGetTsFieldName(): 1},
			Options: options.Index().
				SetExpireAfterSeconds(int32(ttlSec)).
				SetName(idxNameRawEventsTTL),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AppendRawEvents appends events to the raw events collection
// Errors:
// - oerrs.ErrInternal: if insertion failed or nif amount of inserted events is less than input.
func (r *Repo) AppendRawEvents(ctx context.Context, events []Event) error {
	v, err := r.client.
		DB().
		Collection(collectionNameRawEvents).
		InsertMany(ctx, events)
	if err != nil {
		return oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
	}

	if len(v.InsertedIDs) != len(events) {
		return oerrs.NewTErrf(
			ctx,
			"inserted %d events, expected %d; %w",
			len(v.InsertedIDs),
			len(events),
			oerrs.ErrInternal,
		)
	}

	return nil
}
