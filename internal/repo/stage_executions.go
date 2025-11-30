package repo

import (
	"context"
	"errors"
	"os"
	"strconv"

	oerrs "github.com/LastSprint/pipetank/pkg/observability/errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	collectionNameSingleStageExec = "single_stage_exec"
	idxNameSingleStageTTL         = "ttl_stage_exec"
)

func getSingleStageExecutionFilter(
	id, executorID, stageName string,
	processID uuid.UUID,
) bson.M {
	return bson.M{
		SingleStageExecutionEventExecutorIDFieldName():   executorID,
		SingleStageExecutionEventRawStageNameFieldName(): stageName,
		SingleStageExecutionEventProcessIDFieldName():    processID,
		"_id": id,
	}
}

func (r *Repo) createSingleStageIndexes(ctx context.Context) error {
	ttlSec, err := strconv.ParseInt(os.Getenv("STAGE_EXECUTIONS_TTL_SECONDS"), 10, 32)
	if err != nil {
		return oerrs.NewTErr(ctx, err, oerrs.ErrBadInput)
	}

	err = r.client.CreateOrUpdateTTLIndex(
		ctx,
		collectionNameSingleStageExec,
		idxNameSingleStageTTL,
		int32(ttlSec),
		mongo.IndexModel{
			Keys: bson.M{SingleStageExecutionEventUpdateAtFieldName(): 1},
			Options: options.Index().
				SetExpireAfterSeconds(int32(ttlSec)).
				SetName(idxNameSingleStageTTL),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) StartSingleStageExecution(
	ctx context.Context,
	event SingleStageExecutionEvent,
) error {
	if event.UpdatedAt.IsZero() {
		event.UpdatedAt = r.clock()
	}

	v, err := r.client.
		DB().
		Collection(collectionNameSingleStageExec).
		UpdateOne(
			ctx,
			event,
			options.
				UpdateOne().
				SetUpsert(true),
		)
	if err != nil {
		return oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
	}

	if v.ModifiedCount != 1 || v.UpsertedCount != 1 {
		return oerrs.NewTErrf(ctx, "unexpected result of upsert: %v", v)
	}

	return nil
}

func (r *Repo) UpdateSingleStageExecution(
	ctx context.Context,
	id, executorID, stageName string,
	processID uuid.UUID,
	events []Event,
) error {
	updateOperation := bson.M{
		"$set":      bson.M{SingleStageExecutionEventUpdateAtFieldName(): r.clock()},
		"$addToSet": bson.M{SingleStageExecutionEventUpdatesFieldName(): bson.M{"$each": events}},
	}

	v, err := r.client.
		DB().
		Collection(collectionNameSingleStageExec).
		UpdateOne(ctx, getSingleStageExecutionFilter(id, executorID, stageName, processID), updateOperation)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return oerrs.NewTErrf(ctx, "no such execution: %w", oerrs.ErrNotFound)
	}

	if err != nil {
		return oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
	}

	if v.MatchedCount != 1 {
		return oerrs.NewTErrf(
			ctx,
			"expected to update 1 document, got %d:%w",
			v.MatchedCount,
			oerrs.ErrNotFound,
		)
	}

	if v.ModifiedCount != 1 {
		return oerrs.NewTErrf(ctx, "expected to update 1 document, got %d", v.ModifiedCount)
	}

	if v.UpsertedCount != 0 {
		return oerrs.NewTErrf(ctx, "expected to upsert 0 documents, got %d", v.UpsertedCount)
	}

	return err
}

func (r *Repo) FinishSingleStageExecution(
	ctx context.Context,
	event Event,
	isSuccess bool,
) error {
	filter := getSingleStageExecutionFilter(
		event.ExecID,
		event.ExecutorID,
		event.Stage.Name,
		event.ProcessID,
	)
	update := bson.M{
		"$set": bson.M{
			SingleStageExecutionEventUpdateAtFieldName():   r.clock(),
			SingleStageExecutionEventEndFieldName():        event,
			SingleStageExecutionEventIsFinishedFieldName(): true,
			SingleStageExecutionEventIsSuccessFieldName():  isSuccess,
		},
	}

	v, err := r.client.DB().Collection(collectionNameSingleStageExec).UpdateOne(ctx, filter, update)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return oerrs.NewTErrf(ctx, "no such execution: %w", oerrs.ErrNotFound)
	}

	if err != nil {
		return oerrs.NewTErr(ctx, err, oerrs.ErrInternal)
	}

	if v.MatchedCount != 1 {
		return oerrs.NewTErrf(ctx, "expected to update 1 document, got %d", v.MatchedCount)
	}

	return nil
}
