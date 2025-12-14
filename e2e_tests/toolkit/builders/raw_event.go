//go:build test

package builders

import (
	"fmt"
	"testing"
	"time"

	"github.com/LastSprint/pipetank/internal/repo"
	"github.com/LastSprint/pipetank/pkg/mdb"
	"github.com/stretchr/testify/require"
)

type RawEventBuilder struct {
	Event repo.Event

	t *testing.T
}

func NewRawEventBuilder(t *testing.T) *RawEventBuilder {
	return &RawEventBuilder{
		t: t,
	}
}

func (b *RawEventBuilder) WithWorkerID(v string) *RawEventBuilder {
	b.Event.WorkerID = v
	return b
}

func (b *RawEventBuilder) WithStageExecutionID(v string) *RawEventBuilder {
	b.Event.StageExecutionID = v
	return b
}

func (b *RawEventBuilder) WithProcessID(v string) *RawEventBuilder {
	b.Event.ProcessID = v
	return b
}

func (b *RawEventBuilder) WithExecutionID(v string) *RawEventBuilder {
	b.Event.ExecutionID = v
	return b
}

func (b *RawEventBuilder) WithPlainExecutionID(v string) *RawEventBuilder {
	b.Event.ExecutionID = fmt.Sprintf("%s_%s_%s", b.Event.ProcessID, b.Event.WorkerID, v)
	return b
}

func (b *RawEventBuilder) WithStage(v repo.RawStage) *RawEventBuilder {
	b.Event.Stage = v
	return b
}

func (b *RawEventBuilder) WithTs(v time.Time) *RawEventBuilder {
	b.Event.Ts = v
	return b
}

func (b *RawEventBuilder) WithKind(v repo.EventKind) *RawEventBuilder {
	b.Event.Kind = v
	return b
}

func (b *RawEventBuilder) WithKindStarted() *RawEventBuilder {
	b.Event.Kind = repo.EventKindStageStarted
	return b
}

func (b *RawEventBuilder) WithKindFinished() *RawEventBuilder {
	b.Event.Kind = repo.EventKindStageFinished
	return b
}

func (b *RawEventBuilder) WithKindUpdate() *RawEventBuilder {
	b.Event.Kind = repo.EventKindGenericUpdate
	return b
}

func (b *RawEventBuilder) WithStatus(v repo.EventStatus) *RawEventBuilder {
	b.Event.Status = v
	return b
}

func (b *RawEventBuilder) WithStatusSuccess() *RawEventBuilder {
	b.Event.Status = repo.EventStatusSuccess
	return b
}

func (b *RawEventBuilder) WithStatusFailure() *RawEventBuilder {
	b.Event.Status = repo.EventStatusFailure
	return b
}

func (b *RawEventBuilder) WithInput() *RawEventBuilder {
	inputData := map[string]any{
		"type":               "input",
		"execution_id":       b.Event.ExecutionID,
		"process_id":         b.Event.ProcessID,
		"stage_execution_id": b.Event.StageExecutionID,
		"worker_id":          b.Event.WorkerID,
		"kind":               b.Event.Kind,
	}

	bts, err := mdb.MarshalBson(inputData)
	if err != nil {
		require.NoError(b.t, err)
	}

	b.Event.Input = bts
	return b
}

func (b *RawEventBuilder) WithOutput() *RawEventBuilder {
	toJSON := map[string]any{
		"type":               "output",
		"execution_id":       b.Event.ExecutionID,
		"process_id":         b.Event.ProcessID,
		"stage_execution_id": b.Event.StageExecutionID,
		"worker_id":          b.Event.WorkerID,
		"kind":               b.Event.Kind,
	}

	bts, err := mdb.MarshalBson(toJSON)
	if err != nil {
		require.NoError(b.t, err)
	}

	b.Event.Output = bts
	return b
}

func (b *RawEventBuilder) WithFailure() *RawEventBuilder {
	b.Event.Kind = repo.EventKindStageFinished
	b.Event.Status = repo.EventStatusFailure
	toJSON := map[string]any{
		"type":               "failure",
		"execution_id":       b.Event.ExecutionID,
		"process_id":         b.Event.ProcessID,
		"stage_execution_id": b.Event.StageExecutionID,
		"worker_id":          b.Event.WorkerID,
		"kind":               b.Event.Kind,
	}
	bts, err := mdb.MarshalBson(toJSON)
	if err != nil {
		require.NoError(b.t, err)
	}

	b.Event.Failure = bts
	return b
}

func (b *RawEventBuilder) WithMetadata() *RawEventBuilder {
	toJSON := map[string]any{
		"type":               "metadata",
		"execution_id":       b.Event.ExecutionID,
		"process_id":         b.Event.ProcessID,
		"stage_execution_id": b.Event.StageExecutionID,
		"worker_id":          b.Event.WorkerID,
		"kind":               b.Event.Kind,
	}

	bts, err := mdb.MarshalBson(toJSON)
	if err != nil {
		require.NoError(b.t, err)
	}

	b.Event.Metadata = bts
	return b
}

/// ------- Stage

func (b *RawEventBuilder) WithStageName(v string) *RawEventBuilder {
	b.Event.Stage.Name = v
	return b
}

func (b *RawEventBuilder) WithStageDescription(v string) *RawEventBuilder {
	b.Event.Stage.Description = v
	return b
}

// ------- Helpers

func (b *RawEventBuilder) Copy() *RawEventBuilder {
	return &RawEventBuilder{
		Event: b.Event,
		t:     b.t,
	}
}

// ------- Getters
