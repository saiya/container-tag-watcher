package watcher

import (
	"context"

	"github.com/saiya/container_tag_watcher/config"
	"github.com/saiya/container_tag_watcher/logger"
)

type Watcher struct {
	rootCtx       context.Context
	rootCtxCloser context.CancelFunc

	callback ImageUpdateCallback
}

type ImageUpdateCallback func(targetName string)

func NewWatcher(ctx context.Context, cfg *config.Config, callback ImageUpdateCallback) *Watcher {
	rootCtx, rootCtxCloser := context.WithCancel(ctx)
	w := &Watcher{
		rootCtx:       rootCtx,
		rootCtxCloser: rootCtxCloser,

		callback: callback,
	}

	for name, target := range cfg.Targets {
		newWorker(ctx, name, target, callback)
	}
	return w
}

func (w *Watcher) Close() {
	logger.Get().Debugw("Closing watcher...")

	w.rootCtxCloser()
}
