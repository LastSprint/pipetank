//go:build test

package client

import (
	"testing"

	"github.com/LastSprint/pipetank/pkg/client/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TestClient struct {
	t         *testing.T
	apiClient proto.APIClient

	eventsStream grpc.ClientStreamingClient[proto.ClientCommand, emptypb.Empty]
}

func NewTestClient(t *testing.T, serverAddr string) *TestClient {
	t.Helper()
	cn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	require.NoError(t, err)
	cl := proto.NewAPIClient(cn)
	stream, err := cl.Stream(t.Context())
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = stream.CloseSend()
		_ = cn.Close()
	})

	return &TestClient{
		t:            t,
		apiClient:    cl,
		eventsStream: stream,
	}
}

func (c *TestClient) SendRawEvent(t *testing.T, events ...*proto.RawEvent) {
	t.Helper()
	batch := &proto.RawEvents{
		Events: events,
	}

	val := proto.ClientCommand{Cmd: &proto.ClientCommand_Events{Events: batch}}

	err := c.eventsStream.Send(&val)
	require.NoError(t, err)
}
