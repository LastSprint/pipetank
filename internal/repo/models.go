package repo

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	EventKind   int
	EventStatus int
)

const (
	// EventKindStageStarted describes beginning of a stage execution.
	EventKindStageStarted EventKind = 0
	// EventKindStageFinished describes end of a stage execution.
	EventKindStageFinished EventKind = 1
	// EventKindGenericUpdate describes any other event that might happen during the execution.
	EventKindGenericUpdate EventKind = 2

	EventStatusUnknown EventStatus = 0
	// EventStatusSuccess describes successful stage execution.
	EventStatusSuccess EventStatus = 1
	// EventStatusFailure describes failed stage execution
	// If a stage has failed, the output field will be empty.
	EventStatusFailure EventStatus = 1
)

func (kind EventKind) IsValid() bool {
	return kind >= EventKindStageStarted && kind <= EventKindGenericUpdate
}

func (status EventStatus) IsValid() bool {
	return status >= EventStatusSuccess && status <= EventStatusFailure
}

type RawStage struct {
	Name        string `bson:"n"`
	Description string `bson:"d,omitempty"`
}

func (s RawStage) Validate() error {
	var resultErr error

	if len(s.Name) == 0 {
		resultErr = errors.Join(resultErr, errors.New("Stage.Name must be set"))
	}

	return resultErr
}

func RawStageNameFieldName() string {
	return "n"
}

type Event struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	ProcessID        string `bson:"pid"`
	ExecutionID      string `bson:"eid"`
	StageExecutionID string `bson:"seid"`

	WorkerID string `bson:"wid"`

	Stage  RawStage    `bson:"s"`
	Ts     time.Time   `bson:"ts"`
	Kind   EventKind   `bson:"k"`
	Status EventStatus `bson:"st"`

	// Input is optional and may be set only for EventKindStageStarted
	Input bson.Raw `bson:"i,omitempty"`
	// Output is optional and may be set only for EventKindStageFinished
	Output bson.Raw `bson:"o,omitempty"`
	// Failure is optional and may be set only for EventKindStageFinished disregard the status
	Failure bson.Raw `bson:"f,omitempty"`
	// Metadata is an optional field that can be used to store any additional information.
	Metadata bson.Raw `bson:"m,omitempty"`
}

func (e Event) Validate() error {
	var resultErr error

	if !e.ID.IsZero() {
		resultErr = errors.Join(resultErr, errors.New("ID must be empty"))
	}

	if len(e.ProcessID) == 0 {
		resultErr = errors.Join(resultErr, errors.New("ProcessID must be set"))
	}

	if len(e.WorkerID) == 0 {
		resultErr = errors.Join(resultErr, errors.New("WorkerID must be set"))
	}

	if len(e.ExecutionID) == 0 {
		resultErr = errors.Join(resultErr, errors.New("ExecutionID must be set"))
	}

	if len(e.StageExecutionID) == 0 {
		resultErr = errors.Join(resultErr, errors.New("StageExecutionID must be set"))
	}

	if e.Ts.IsZero() {
		resultErr = errors.Join(resultErr, errors.New("TS must be set"))
	}

	if !e.Kind.IsValid() {
		resultErr = errors.Join(resultErr, errors.New(".Kind must be valid"))
	}

	if !e.Status.IsValid() {
		resultErr = errors.Join(resultErr, errors.New(".Status must be valid"))
	}

	return errors.Join(resultErr, e.Stage.Validate())
}

func (e Event) Copy() Event {
	var inputCp bson.Raw
	var outputCp bson.Raw
	var failureCp bson.Raw
	var metadataCp bson.Raw

	if len(e.Input) > 0 {
		inputCp = make(bson.Raw, len(e.Input))
		copy(inputCp, e.Input)
	}
	if len(e.Output) > 0 {
		outputCp = make(bson.Raw, len(e.Output))
		copy(outputCp, e.Output)
	}
	if len(e.Failure) > 0 {
		failureCp = make(bson.Raw, len(e.Failure))
		copy(failureCp, e.Failure)
	}
	if len(e.Metadata) > 0 {
		metadataCp = make(bson.Raw, len(e.Metadata))
		copy(metadataCp, e.Metadata)
	}

	return Event{
		ID:               e.ID,
		ProcessID:        e.ProcessID,
		WorkerID:         e.WorkerID,
		ExecutionID:      e.ExecutionID,
		StageExecutionID: e.StageExecutionID,
		Stage:            e.Stage,
		Ts:               e.Ts,
		Kind:             e.Kind,
		Status:           e.Status,
		Input:            inputCp,
		Output:           outputCp,
		Failure:          failureCp,
		Metadata:         metadataCp,
	}
}

func RawEventGetTsFieldName() string {
	return "ts"
}

type SingleStageExecutionEvent struct {
	ProcessID        string `bson:"pid"`
	WorkerID         string `bson:"wid"`
	ExecutionID      string `bson:"eid"`
	StageExecutionID string `bson:"seid"`

	RawStage RawStage `bson:"rs,omitempty"`

	Start   Event   `bson:"s"`
	Updates []Event `bson:"u,omitempty"`
	End     Event   `bson:"e,omitempty"`

	IsFinished bool `bson:"if"`
	IsSuccess  bool `bson:"is"`

	UpdatedAt time.Time `bson:"ua"`
}

func SingleStageExecutionEventUpdateAtFieldName() string {
	return "ua"
}

func SingleStageExecutionEventProcessIDFieldName() string {
	return "oid"
}

func SingleStageExecutionEventExecutionIDFieldName() string {
	return "eid"
}

func SingleStageExecutionEventStageExecutionIDFieldName() string {
	return "seid"
}

func SingleStageExecutionEventIsFinishedFieldName() string {
	return "if"
}

func SingleStageExecutionEventIsSuccessFieldName() string {
	return "is"
}

func SingleStageExecutionEventUpdatesFieldName() string {
	return "u"
}

func SingleStageExecutionEventEndFieldName() string {
	return "e"
}

func SingleStageExecutionEventRawStageNameFieldName() string {
	return "rs" + "." + RawStageNameFieldName()
}
