package main

import (
	"context"

	app "github.com/LastSprint/pipetank/internal/apps/mongodb_change_stream_consumer"
)

var Version = "0.0.1"

func main() {
	ctx := context.Background()
	err := app.Run(ctx)
	if err != nil {
		panic(err)
	}
}
