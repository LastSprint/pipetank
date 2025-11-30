package mdbset

import (
	"fmt"
	"testing"
)

type TestDB struct {
	Port   int
	DBName string
}

func (t TestDB) GetDSN() string {
	return fmt.Sprintf("mongodb://localhost:%d", t.Port)
}

func StartMongoDB(t *testing.T) (tdb TestDB) {
	t.Helper()
	return tdb
}
