package watcher

import (
	"context"
	"time"

	"github.com/saiya/container_tag_watcher/config"
)

type worker struct {
	targetName string
}

func newWorker(ctx context.Context, targetName string, cfg *config.WatchTarget, callback ImageUpdateCallback) *worker {
	w := &worker{targetName: targetName}

	// FIXME: Implement real logic
	// This is dummy logic
	pollingInterval, err := time.ParseDuration(cfg.PollingInterval)
	if err != nil {
		panic(err) // Config validator should reject such config
	}
	go func() {
		ticker := time.NewTicker(pollingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				callback(targetName)
			}
		}
	}()

	return w
}
