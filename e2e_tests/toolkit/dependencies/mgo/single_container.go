//go:build test

package mgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func RunSingleContainer(t *testing.T) string {
	t.Helper()
	container, err := mongodb.Run(t.Context(), "mongo:8")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	dsn, err := container.ConnectionString(t.Context())
	require.NoError(t, err)

	return dsn
}
