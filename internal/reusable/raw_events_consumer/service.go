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
)

type store interface {
	StartSingleStageExecution(
		ctx context.Context,
		event repo.SingleStageExecutionEvent,
	) error

	UpdateSingleStageExecution(
		ctx context.Context,
		processID, executionID, stageExecutionID string,
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

// HandleEvents does all the necessary processing.
//
// !!!IMPORTANT!!! DOES NOT VALIDATE EVENTS.
// You have to use repo.Event.Validate() before calling this function.
//
// It is done this way to save some CPU cycles, improve cache and avoid unnecessary allocations
// (because you have to map events from your source to repo.Event anyway).
func (s *Service) HandleEvents(
	ctx context.Context,
	events []repo.Event,
) error {
	if len(events) == 0 {
		return nil
	}

	ne := normalizeEvents(events)

	var errs error

	for processID, executions := range ne {
		for executionID, stages := range executions {
			for _, stageExecutions := range stages {
				for stageExecutionID, events := range stageExecutions {
					err := s.actOnEvents(ctx, events, processID, executionID, stageExecutionID)
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
	processID processID,
	executionID executionID,
	stageExecutionID stageExecutionID,
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
					ProcessID:        event.ProcessID,
					ExecutionID:      event.ExecutionID,
					StageExecutionID: event.StageExecutionID,
					RawStage:         event.Stage,
					Start:            event,
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
			processID, executionID, stageExecutionID,
			updatesToSave,
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
	processID        = string
	executionID      = string
	stageExecutionID = string
	stageName        = string

	normalizedEvents map[processID]map[executionID]map[stageName]map[stageExecutionID][]repo.Event
)

func normalizeEvents(events []repo.Event) normalizedEvents {
	result := normalizedEvents{}

	for _, event := range events {
		process, ok := result[event.ProcessID]

		if !ok {
			process = map[executionID]map[stageName]map[stageExecutionID][]repo.Event{}
		}

		execution, ok := process[event.ExecutionID]
		if !ok {
			execution = map[stageName]map[stageExecutionID][]repo.Event{}
		}

		stage, ok := execution[event.Stage.Name]
		if !ok {
			stage = map[stageExecutionID][]repo.Event{}
		}

		stage[event.StageExecutionID] = append(stage[event.StageExecutionID], event)
		execution[event.Stage.Name] = stage
		process[event.ExecutionID] = execution
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
