//go:build test

package utils

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func RandomOpenPort(t *testing.T) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	addr := l.Addr().(*net.TCPAddr)

	require.NoError(t, l.Close())

	return addr.Port
}
