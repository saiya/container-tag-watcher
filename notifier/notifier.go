package notifier

import (
	"context"

	"github.com/saiya/container_tag_watcher/config"
	"github.com/saiya/container_tag_watcher/logger"
)

type Notifier struct {
	rootCtx       context.Context
	rootCtxCloser context.CancelFunc

	workers map[string]*worker
}

func NewNotifier(ctx context.Context, cfg *config.Config) *Notifier {
	rootCtx, rootCtxCloser := context.WithCancel(ctx)
	n := &Notifier{
		rootCtx:       rootCtx,
		rootCtxCloser: rootCtxCloser,

		workers: map[string]*worker{},
	}
	for targetName, targetCfg := range cfg.Targets {
		n.workers[targetName] = newWorker(n.rootCtx, targetName, targetCfg)
	}
	return n
}

func (n *Notifier) Close() {
	logger.Get().Debugw("Closing notifier...")
	n.rootCtxCloser()
}

func (n *Notifier) OnImageUpdated(targetName string) {
	n.workers[targetName].OnImageUpdated()
}
