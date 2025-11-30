package main

import (
	"context"

	app "github.com/LastSprint/pipetank/internal/apps/grpc_api"
)

var Version = "0.0.1"

func main() {
	ctx := context.Background()
	err := app.Run(ctx)
	if err != nil {
		panic(err)
	}
}
