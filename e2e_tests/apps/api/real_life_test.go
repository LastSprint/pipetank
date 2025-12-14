//go:build test

package api

import (
	"testing"
	"time"

	"github.com/LastSprint/pipetank/e2e_tests/toolkit/builders"
	"github.com/LastSprint/pipetank/internal/repo"
)

func TestRealLife(t *testing.T) {
	var (
		process1 = "p1"
		process2 = "p2"

		worker1 = "w1"
		worker2 = "w2"
		worker3 = "w3"
	)

	testEnv := RunTestEnv(t, worker1, worker2, worker3)

	stage1 := repo.RawStage{Name: "stage_1", Description: "stage 1 description"}
	stage2 := repo.RawStage{Name: "stage_2", Description: "stage 1 description"}
	stage3 := repo.RawStage{Name: "stage_3", Description: "stage 1 description"}

	beginningOfTime := time.Now()

	// # Client 1 events
	tsStart := beginningOfTime.Add(time.Second * 1)
	step := time.Second * 5

	// ## First Execution | First Process

	reps := []*builders.RawEventBuilder{
		// <editor-fold desc="Client 1 - Process 1 - Execution 1">
		builders.NewRawEventBuilder(t).
			WithTs(tsStart).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 2)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 3)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 4)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithKindFinished().
			WithStatusSuccess().
			WithOutput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 5)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindFinished().
			WithStatusSuccess().
			WithOutput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 6)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindFinished().
			WithStatusSuccess().
			WithOutput(),

		// </editor-fold>

		//<editor-fold desc="Client 1 - Process 1 - Execution 2">
		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 7)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 8)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 9)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 10)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithFailure(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 11)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithFailure(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 12)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithFailure(),

		//</editor-fold>

		//<editor-fold desc="Client 1 - Process 1 - Execution 3">
		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 8)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_3").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 9)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_3").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindUpdate().
			WithMetadata(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 10)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_3").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithFailure().
			WithOutput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 8)).
			WithWorkerID(worker1).
			WithProcessID(process1).
			WithPlainExecutionID("e_3").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindFinished().
			WithStatusSuccess().
			WithOutput(),

		//</editor-fold>

		//<editor-fold desc="Client 2 - Process 1 - Execution 1">
		builders.NewRawEventBuilder(t).
			WithTs(tsStart).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 2)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 3)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindUpdate().
			WithMetadata(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 4)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindUpdate().
			WithMetadata(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 5)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 6)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithKindFinished().
			WithStatusSuccess().
			WithOutput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 7)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindFinished().
			WithStatusSuccess().
			WithOutput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 8)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindFinished().
			WithStatusSuccess().
			WithOutput(),

		//</editor-fold>

		//<editor-fold desc="Client 2 - Process 1 - Execution 2">
		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 9)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 10)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 11)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindFinished().
			WithFailure(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 12)).
			WithWorkerID(worker2).
			WithProcessID(process1).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindFinished().
			WithFailure(),

		//</editor-fold>

		//<editor-fold desc="Client 3 - Process 2 - Execution 1">
		builders.NewRawEventBuilder(t).
			WithTs(tsStart).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 2)).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 4)).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 6)).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage3.Name).
			WithStage(stage3).
			WithFailure(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 8)).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage2.Name).
			WithStage(stage2).
			WithKindFinished().
			WithStatusSuccess().
			WithOutput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 10)).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_1").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindFinished().
			WithOutput(),

		//</editor-fold>

		//<editor-fold desc="Client 3 - Process 2 - Execution 2">
		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 12)).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_2").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindStarted().
			WithInput(),

		//</editor-fold>

		//<editor-fold desc="Client 3 - Process 2 - Execution 3">
		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 14)).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_3").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithKindStarted().
			WithInput(),

		builders.NewRawEventBuilder(t).
			WithTs(tsStart.Add(step * 16)).
			WithWorkerID(worker3).
			WithProcessID(process2).
			WithPlainExecutionID("e_3").
			WithStageExecutionID(stage1.Name).
			WithStage(stage1).
			WithFailure(),

		//</editor-fold>
	}

	_ = testEnv
	_ = reps
}
