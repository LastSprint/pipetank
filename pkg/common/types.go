package common

import "context"

type CallbackFailable[T any] func(ctx context.Context, event T) error
