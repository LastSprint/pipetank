//go:build test

package api

import (
	"testing"

	"github.com/LastSprint/pipetank/e2e_tests/toolkit/dependencies/client"
	testGrpcAPI "github.com/LastSprint/pipetank/e2e_tests/toolkit/dependencies/grpc_api"
	mgo2 "github.com/LastSprint/pipetank/e2e_tests/toolkit/dependencies/mgo"
	"github.com/LastSprint/pipetank/pkg/mdb"
)

type TestEnv struct {
	MdbClient  *mdb.Client
	Clients    map[string]*client.TestClient
	ServerAddr string
}

func RunTestEnv(t *testing.T, clients ...string) *TestEnv {
	dbName := "e2e_test_db"
	mdbDsn := mgo2.RunSingleContainer(t)
	mdbClient := mgo2.InitTestMDBClient(t, mdbDsn, dbName)

	srvAddr := testGrpcAPI.Run(t, mdbDsn, dbName)

	clientsMap := make(map[string]*client.TestClient, len(clients))

	for _, clientKey := range clients {
		clientsMap[clientKey] = client.NewTestClient(t, srvAddr)
	}

	return &TestEnv{
		MdbClient:  mdbClient,
		Clients:    clientsMap,
		ServerAddr: srvAddr,
	}
}
