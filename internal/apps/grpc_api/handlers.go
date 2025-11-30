package grpc_api

import (
	"context"
	"errors"
	"io"

	"github.com/LastSprint/pipetank/internal/repo"
	"github.com/LastSprint/pipetank/pkg/client/proto"
	"github.com/LastSprint/pipetank/pkg/mdb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type eventHandler interface {
	HandleEvents(context.Context, []repo.Event) error
}

type handlers struct {
	proto.UnimplementedAPIServer

	eventHandler eventHandler
}

func newHandlers(
	eventHandler eventHandler,
) *handlers {
	return &handlers{eventHandler: eventHandler}
}

func (h *handlers) SendEvent(
	stream grpc.ClientStreamingServer[proto.RawEvents, emptypb.Empty],
) error {
	for {
		batch, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// the client closed the stream
			break
		}

		if err != nil {
			return status.Errorf(codes.Internal, "recv error: %v", err)
		}

		converted := make([]repo.Event, 0, len(batch.Events))

		err = convertEventsRaw(batch, converted)
		if err != nil {
			return err
		}

		err = h.eventHandler.HandleEvents(stream.Context(), converted)
		if err != nil {
			return err
		}
	}

	if err := stream.SendAndClose(&emptypb.Empty{}); err != nil {
		return status.Errorf(codes.Internal, "send and close error: %v", err)
	}

	return nil
}

func convertEventsRaw(batch *proto.RawEvents, converted []repo.Event) error {
	for _, ev := range batch.Events {
		bsonInput, err := mdb.JSONtoBSON(ev.Input)
		if err != nil {
			return status.Errorf(codes.Internal, "bson conversion error: %v", err)
		}
		bsonOutput, err := mdb.JSONtoBSON(ev.Output)
		if err != nil {
			return status.Errorf(codes.Internal, "bson conversion error: %v", err)
		}
		bsonFailure, err := mdb.JSONtoBSON(ev.Failure)
		if err != nil {
			return status.Errorf(codes.Internal, "bson conversion error: %v", err)
		}
		bsonMetadata, err := mdb.JSONtoBSON(ev.Metadata)
		if err != nil {
			return status.Errorf(codes.Internal, "bson conversion error: %v", err)
		}

		it := repo.Event{
			ExecID:     ev.GetExecID(),
			ProcessID:  uuid.UUID(ev.ProcessID),
			ExecutorID: ev.GetExecutorID(),
			Stage: repo.RawStage{
				Name:        ev.GetStage().GetName(),
				Description: ev.GetStage().GetDescription(),
			},
			Ts:       ev.GetTs().AsTime(),
			Kind:     repo.EventKind(ev.GetKind()),
			Status:   repo.EventStatus(ev.GetStatus()),
			Input:    bsonInput,
			Output:   bsonOutput,
			Failure:  bsonFailure,
			Metadata: bsonMetadata,
		}

		converted = append(converted, it)
	}

	return nil
}
