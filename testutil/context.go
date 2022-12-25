package testutil

import (
	"context"
	"testing"
)

func Context(t *testing.T) context.Context {
	ctx, close := context.WithCancel(context.Background())
	t.Cleanup(close)
	return ctx
}
