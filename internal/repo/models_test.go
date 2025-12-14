package repo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestEventValidate(t *testing.T) {
	validEvent := Event{
		ExecID:     "exec-1",
		ProcessID:  uuid.New(),
		ExecutorID: "worker-42",
		Stage: RawStage{
			Name: "stage-name",
		},
		Ts:     time.Now(),
		Kind:   EventKindStageStarted,
		Status: EventStatusSuccess,
	}

	testCases := []struct {
		name     string
		mutate   func(Event) Event
		errParts []string
	}{
		{
			name:   "valid event",
			mutate: func(e Event) Event { return e },
		},
		{
			name: "non-empty id",
			mutate: func(e Event) Event {
				e.ID = bson.NewObjectID()
				return e
			},
			errParts: []string{"ID must be empty"},
		},
		{
			name: "empty exec id",
			mutate: func(e Event) Event {
				e.ExecID = ""
				return e
			},
			errParts: []string{"ExecID must be set"},
		},
		{
			name: "nil process id",
			mutate: func(e Event) Event {
				e.ProcessID = uuid.Nil
				return e
			},
			errParts: []string{"ProcessID must be set"},
		},
		{
			name: "empty executor id",
			mutate: func(e Event) Event {
				e.ExecutorID = ""
				return e
			},
			errParts: []string{"ExecutorID must be set"},
		},
		{
			name: "empty stage name",
			mutate: func(e Event) Event {
				e.Stage.Name = ""
				return e
			},
			errParts: []string{"Stage.Name must be set"},
		},
		{
			name: "zero timestamp",
			mutate: func(e Event) Event {
				e.Ts = time.Time{}
				return e
			},
			errParts: []string{"Ts must be set"},
		},
		{
			name: "invalid kind",
			mutate: func(e Event) Event {
				e.Kind = EventKind(100)
				return e
			},
			errParts: []string{"Kind must be valid"},
		},
		{
			name: "invalid status",
			mutate: func(e Event) Event {
				e.Status = EventStatus(-1)
				return e
			},
			errParts: []string{"Status must be valid"},
		},
		{
			name: "multiple validation errors",
			mutate: func(e Event) Event {
				e.ExecID = ""
				e.Kind = EventKind(10)
				e.Status = EventStatus(-5)
				return e
			},
			errParts: []string{
				"ExecID must be set",
				"Kind must be valid",
				"Status must be valid",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			event := tc.mutate(validEvent)

			err := event.Validate()

			if len(tc.errParts) == 0 {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)

			for _, expected := range tc.errParts {
				assert.ErrorContains(t, err, expected)
			}
		})
	}
}
