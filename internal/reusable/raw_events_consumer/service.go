// Package raweventsconsumer (Package Raw Events Consumer)
// Encapsulates logic for handle raw events from any generic source.
//
// It provides only one function - HandleEvents which accept any kind of repo.Event in any order and any size
// and then perform all necessary transformation and saving.
package raweventsconsumer

import (
	"context"
	"errors"
	"slices"

	"github.com/LastSprint/pipetank/internal/repo"
	"github.com/google/uuid"
)

type store interface {
	StartSingleStageExecution(
		ctx context.Context,
		event repo.SingleStageExecutionEvent,
	) error

	UpdateSingleStageExecution(
		ctx context.Context,
		id, executorID, stageName string,
		processID uuid.UUID,
		events []repo.Event,
	) error

	FinishSingleStageExecution(
		ctx context.Context,
		event repo.Event,
		isSuccess bool,
	) error
}

type Service struct {
	store store
}

func NewService(store store) *Service {
	return &Service{store: store}
}

func (s *Service) HandleEvents(
	ctx context.Context,
	events []repo.Event,
) error {
	if len(events) == 0 {
		return nil
	}

	ne := normalizeEvents(events)

	var errs error

	for processID, executors := range ne {
		for executorID, stages := range executors {
			for stageName, executions := range stages {
				for execID, events := range executions {
					err := s.actOnEvents(ctx, events, execID, executorID, stageName, processID)
					if err != nil {
						errs = errors.Join(errs, err)
					}
				}
			}
		}
	}

	return errs
}

func (s *Service) actOnEvents(
	ctx context.Context,
	events []repo.Event,
	execID executionID,
	executorID executorID,
	stageName stageName,
	processID processID,
) error {
	updatesToSave := make([]repo.Event, 0)
	needToEndEvent := false
	var endingEvent repo.Event

	var errs error

	for _, event := range events {
		switch event.Kind {
		case repo.EventKindStageStarted:
			err := s.store.StartSingleStageExecution(
				ctx,
				repo.SingleStageExecutionEvent{
					ID:         event.ExecID,
					ProcessID:  event.ProcessID,
					ExecutorID: event.ExecutorID,
					RawStage:   event.Stage,
					Start:      event,
				},
			)
			if err != nil {
				errs = errors.Join(errs, err)
			}

		case repo.EventKindStageFinished:
			endingEvent = event
			needToEndEvent = true
		case repo.EventKindGenericUpdate:
			updatesToSave = append(updatesToSave, event)
		}
	}

	if len(updatesToSave) > 0 {
		err := s.store.UpdateSingleStageExecution(
			ctx,
			execID, executorID, stageName, processID, updatesToSave,
		)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if needToEndEvent {
		err := s.store.FinishSingleStageExecution(
			ctx,
			endingEvent,
			endingEvent.Kind == repo.EventKindStageFinished,
		)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

type (
	processID   = uuid.UUID
	executorID  = string
	executionID = string
	stageName   = string

	normalizedEvents map[processID]map[executorID]map[stageName]map[executionID][]repo.Event
)

func normalizeEvents(events []repo.Event) normalizedEvents {
	result := normalizedEvents{}

	for _, event := range events {
		process, ok := result[event.ProcessID]

		if !ok {
			process = map[executorID]map[stageName]map[executionID][]repo.Event{}
		}

		executor, ok := process[event.ExecutorID]
		if !ok {
			executor = map[stageName]map[executionID][]repo.Event{}
		}

		stage, ok := executor[event.Stage.Name]
		if !ok {
			stage = map[executionID][]repo.Event{}
		}

		stage[event.ExecID] = append(stage[event.ExecID], event)
		executor[event.Stage.Name] = stage
		process[event.ExecutorID] = executor
		result[event.ProcessID] = process
	}

	sortNormalizedEvents(result)

	return result
}

func sortNormalizedEvents(ne normalizedEvents) {
	for _, executors := range ne {
		for _, stages := range executors {
			for _, executions := range stages {
				for _, events := range executions {
					slices.SortFunc(events, func(a, b repo.Event) int {
						return a.Ts.Compare(b.Ts)
					})
				}
			}
		}
	}
}
