package repo

import (
	"time"

	"github.com/google/uuid"
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

	// EventStatusSuccess describes successful stage execution.
	EventStatusSuccess EventStatus = 0
	// EventStatusFailure describes failed stage execution
	// If a stage has failed, the output field will be empty.
	EventStatusFailure EventStatus = 1
)

type RawStage struct {
	Name        string `bson:"n"`
	Description string `bson:"d,omitempty"`
}

func RawStageNameFieldName() string {
	return "n"
}

type Event struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	// ExecID is a unique identifier of the execution (from EventKindStageStarted till EventKindStageFinished)
	// It's like lifetime ID of a single execution of a single stage
	ExecID     string      `bson:"exid"`
	ProcessID  uuid.UUID   `bson:"pid"`
	ExecutorID string      `bson:"eid"`
	Stage      RawStage    `bson:"s"`
	Ts         time.Time   `bson:"ts"`
	Kind       EventKind   `bson:"k"`
	Status     EventStatus `bson:"st"`

	// Input is optional and may be set only for EventKindStageStarted
	Input bson.Raw `bson:"i,omitempty"`
	// Output is optional and may be set only for EventKindStageFinished
	Output bson.Raw `bson:"o,omitempty"`
	// Failure is optional and may be set only for EventKindStageFinished disregard the status
	Failure bson.Raw `bson:"f,omitempty"`
	// Metadata is an optional field that can be used to store any additional information.
	Metadata bson.Raw `bson:"m,omitempty"`
}

func RawEventGetTsFieldName() string {
	return "ts"
}

type Stage struct {
	ID          string `bson:"_id"`
	Name        string `bson:"n"`
	Description string `bson:"d,omitempty"`
}

type SingleStageExecutionEvent struct {
	ID         string    `bson:"_id"`
	ProcessID  uuid.UUID `bson:"pid"`
	ExecutorID string    `bson:"eid"`

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

func SingleStageExecutionEventExecutorIDFieldName() string {
	return "eid"
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
