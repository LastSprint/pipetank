//go:build test

package grpc_api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/LastSprint/pipetank/e2e_tests/toolkit/utils"
	appGrpcAPI "github.com/LastSprint/pipetank/internal/apps/grpc_api"
	"github.com/LastSprint/pipetank/pkg/client/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Run(
	t *testing.T,
	mdbDSN, mdbDBName string,
) string {
	t.Helper()
	port := utils.RandomOpenPort(t)
	srvAddr := fmt.Sprintf("localhost:%d", port)
	go run(t, port, mdbDSN, mdbDBName)

	nc, err := grpc.NewClient(srvAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := proto.NewAPIClient(nc)

	for {
		select {
		case <-t.Context().Done():
			return srvAddr
		case <-time.After(time.Second * 5):
			require.FailNow(t, "server was not ready in time")
		case <-time.After(time.Millisecond * 100):
			_, err := client.HealthCheck(t.Context(), &emptypb.Empty{})
			fmt.Println(err)
			if err == nil {
				return srvAddr
			}
		}
	}
}

func run(t *testing.T, port int, mdbDSN, mdbDBName string) {
	appCfg := appGrpcAPI.Config{Port: port}

	t.Setenv("MDB_DSN", mdbDSN)
	t.Setenv("MDB_DB", mdbDBName)
	t.Setenv("MDB_MAX_CONNECTIONS", "100")
	t.Setenv("APP_NAME", "test-server")
	t.Setenv("MDB_CONNECTION_TIMEOUT", "15s")
	t.Setenv("MDB_MAX_CONNECTING", "100")
	t.Setenv("MDB_MAX_POOL_SIZE", "100")
	t.Setenv("MDB_MIN_POOL_SIZE", "100")
	t.Setenv("STAGE_EXECUTIONS_TTL_SECONDS", strconv.Itoa(60*60*24*30))
	t.Setenv("RAW_EVENTS_TTL_SECONDS", strconv.Itoa(60*60*24*30))

	err := appGrpcAPI.RunWithConfig(t.Context(), appCfg)
	if errors.Is(err, context.Canceled) {
		return
	}
	require.NoError(t, err)
}
