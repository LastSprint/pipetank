//go:build test

package mgo

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/LastSprint/pipetank/pkg/mdb"
	"github.com/stretchr/testify/require"
)

func InitTestMDBClient(t *testing.T, dsn, db string) *mdb.Client {
	t.Helper()
	cfg := mdb.Config{
		DSN:               dsn,
		DB:                db,
		MaxConnections:    100,
		AppName:           "e2e",
		ConnectionTimeout: time.Second * 15,
		MaxConnecting:     100,
		MaxPoolSize:       100,
		MinPoolSize:       100,
	}

	t.Setenv("STAGE_EXECUTIONS_TTL_SECONDS", strconv.Itoa(60*60*24*30))
	t.Setenv("RAW_EVENTS_TTL_SECONDS", strconv.Itoa(60*60*24*30))

	client, err := mdb.NewClientWithConfig(cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = client.Close(context.Background())
	})

	return client
}
