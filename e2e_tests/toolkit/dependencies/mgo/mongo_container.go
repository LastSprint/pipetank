//go:build test

package mgo

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	mongoImage     = "mongo:8"
	replicaSetName = "rs0"
)

type MongoDependency struct {
	DSN string

	Primary    testcontainers.Container
	Secondary1 testcontainers.Container
	Secondary2 testcontainers.Container
	Network    *testcontainers.DockerNetwork
}

// StartRS3 starts a simple RS3 cluster where only the primary is published on a random host port.
// All nodes are attached to the same dedicated docker network to keep traffic isolated.
func StartRS3(t *testing.T) MongoDependency {
	t.Helper()

	dep := MongoDependency{}

	dockerNetwork, err := network.New(
		t.Context(),
		network.WithDriver("bridge"),
	)

	primaryPort := nat.Port("27017")

	require.NoError(t, err)

	dep.Network = dockerNetwork

	primary, _ := startMongoContainer(t, dockerNetwork.Name, &primaryPort)
	dep.Primary = primary

	configureReplicaSet(t, primary, "27017")
	waitForPrimary(t, primary)

	mappedPort, err := primary.MappedPort(t.Context(), primaryPort)
	require.NoError(t, err)

	dep.DSN = fmt.Sprintf(
		"mongodb://%s:%s/?replicaSet=%s",
		"localhost",
		mappedPort.Port(),
		replicaSetName,
	)
	return dep
}

func startMongoContainer(
	t *testing.T,
	networkName string,
	expose *nat.Port,
) (testcontainers.Container, string) {
	t.Helper()

	request := testcontainers.ContainerRequest{
		Image:    mongoImage,
		Cmd:      []string{"mongod", "--replSet", replicaSetName, "--bind_ip_all"},
		Networks: []string{networkName},
		WaitingFor: wait.ForExec(
			[]string{
				"mongosh",
				"--quiet",
				"mongodb://localhost:27017",
				"--eval",
				"db.adminCommand('ping')",
			},
		),
	}

	if expose != nil {
		request.ExposedPorts = []string{expose.Port()}
	}

	container, err := testcontainers.GenericContainer(
		t.Context(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		},
	)
	require.NoError(t, err)

	name, err := container.Name(t.Context())
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = container.Terminate(t.Context())
	})

	return container, strings.TrimPrefix(name, "/")
}

func configureReplicaSet(t *testing.T, primary testcontainers.Container, port string) {
	t.Helper()

	initScript := fmt.Sprintf(
		"rs.initiate({_id:'%s',members:[{_id:0,host:'localhost:27017',priority:2}]})",
		replicaSetName,
	)

	execMongoCommand(t, primary, initScript)
}

func waitForPrimary(t *testing.T, primary testcontainers.Container) {
	t.Helper()

	timeout := time.After(60 * time.Second)
	ctx := t.Context()

	for {
		select {
		case <-ctx.Done():
			require.NoError(t, ctx.Err())
		case <-timeout:
			require.FailNow(t, "mongo primary was not ready in time")
		case <-time.After(100 * time.Millisecond):
			execRes := execMongoCommand(t, primary, "db.hello().isWritablePrimary")
			if strings.Contains(execRes, "true") {
				return
			}
		}
	}
}

func execMongoCommand(t *testing.T, container testcontainers.Container, script string) string {
	t.Helper()
	exitStatus, resultReader, err := container.Exec(
		t.Context(),
		[]string{"mongosh", "--quiet", "--eval", script},
	)
	require.NoError(t, err)

	bt, err := io.ReadAll(resultReader)
	require.NoError(t, err)

	require.Equal(
		t,
		0,
		exitStatus,
		"mongo command exited with non-zero status: %d, output: %s",
		exitStatus,
		string(bt),
	)

	return string(bt)
}
