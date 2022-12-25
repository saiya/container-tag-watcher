package notifier

import (
	"context"

	"github.com/saiya/container_tag_watcher/config"
)

type worker struct {
	name string

	onImageUpdate *WorkerEventHandler
}

func newWorker(ctx context.Context, name string, cfg *config.WatchTarget) *worker {
	w := &worker{
		name: name,

		onImageUpdate: &WorkerEventHandler{
			TargetName: name,
			EventName:  "onImageUpdate",

			BacklogLimit:    *cfg.BacklogLimit,
			Commands:        make([][]string, len(cfg.Commands)),
			ContinueOnError: cfg.ContinueOnError,
		},
	}
	for i, cmdRaw := range cfg.Commands {
		w.onImageUpdate.Commands[i] = cmdRaw.([]string)
	}

	w.onImageUpdate.Start(ctx)
	return w
}

func (w *worker) OnImageUpdated() {
	w.onImageUpdate.Enqueue()
}
