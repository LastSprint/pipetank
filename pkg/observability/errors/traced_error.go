package errors

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	octx "github.com/LastSprint/pipetank/pkg/observability/ctx"
)

const (
	tracedErrorCallersSkip = 3
)

var _tracedErrCapturedSizeDepth = 32

func SetTracedErrCapturedSizeDepth(val int) {
	_tracedErrCapturedSizeDepth = val
}

type TracedError struct {
	nested error

	traceID     string
	stackFrames []frame
}

func NewTErr(ctx context.Context, err error, signals ...error) *TracedError {
	return &TracedError{
		nested:      errors.Join(err, errors.Join(signals...)),
		traceID:     octx.GetTraceID(ctx),
		stackFrames: callers(),
	}
}

func NewTErrf(ctx context.Context, format string, args ...any) *TracedError {
	return &TracedError{
		//nolint:err113
		nested:      fmt.Errorf(format, args...),
		traceID:     octx.GetTraceID(ctx),
		stackFrames: callers(),
	}
}

func (e *TracedError) WithSignal(err error) *TracedError {
	return &TracedError{
		nested:      errors.Join(e.nested, err),
		traceID:     e.traceID,
		stackFrames: e.stackFrames,
	}
}

func (e *TracedError) Unwrap() error {
	return e.nested
}

func (e *TracedError) String() string {
	return fmt.Sprintf("%s\n%s", e.nested.Error(), e.GetFStackTrace())
}

func (e *TracedError) Error() string {
	return fmt.Sprintf("%s\n%s", e.nested.Error(), e.GetFStackTrace())
}

func (e *TracedError) GetTraceID() string {
	return e.traceID
}

func (e *TracedError) GetFStackTrace() string {
	traceStr := ""
	for _, fr := range e.stackFrames {
		traceStr = fmt.Sprintf("%s\t%s\n", traceStr, fr)
	}
	return fmt.Sprintf("%s\n%+v", e.traceID, traceStr)
}

func callers() []frame {
	callers := make([]uintptr, _tracedErrCapturedSizeDepth)
	n := runtime.Callers(tracedErrorCallersSkip, callers)

	if n == 0 {
		return nil
	}

	result := make([]frame, 0, n)

	frames := runtime.CallersFrames(callers)
	for {
		frval, more := frames.Next()

		locFrVal := frame{
			file: frval.File,
			fn:   frval.Function,
			ln:   frval.Line,
		}

		result = append(result, locFrVal)

		if !more {
			break
		}
	}

	return result
}
