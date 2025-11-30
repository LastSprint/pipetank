package repo

import (
	"context"

	"github.com/LastSprint/pipetank/pkg/common"
	"github.com/LastSprint/pipetank/pkg/mdb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RawEventWatchModel struct {
	Record Event
	Token  mdb.ResumeTokenProvider
}

func (r *Repo) WatchRawEvents(
	ctx context.Context,
	action common.CallbackFailable[RawEventWatchModel],
) error {
	return mdb.RunChangeStream(
		ctx,
		r.client,
		collectionNameRawEvents,
		[]bson.M{{"$match": bson.M{"operationType": "insert"}}},
		func(ctx context.Context, token mdb.ResumeTokenProvider, doc Event) error {
			return action(ctx, RawEventWatchModel{Record: doc, Token: token})
		},
	)
}
