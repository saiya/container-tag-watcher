package watcher

import (
	"context"
	"sync"
	"time"

	"github.com/saiya/container_tag_watcher/config"
	cr "github.com/saiya/container_tag_watcher/containerregistryclient"
	"github.com/saiya/container_tag_watcher/logger"
)

type worker struct {
	targetName string

	lock          sync.Mutex
	lastImageHash string
}

func (w *worker) loggerAttrs() []interface{} {
	return []interface{}{"targetName", w.targetName}
}

func (w *worker) updateImageHash(imageHash string) bool {
	w.lock.Lock()
	defer w.lock.Unlock()

	changed := w.lastImageHash != "" && w.lastImageHash != imageHash
	w.lastImageHash = imageHash
	return changed
}

func newWorker(ctx context.Context, targetName string, cfg *config.WatchTarget, crClient cr.Client, callback ImageUpdateCallback) *worker {
	w := &worker{targetName: targetName}

	poll := func() {
		imageHash, err := crClient.GetImageDigest(cfg.Platform, targetName)
		if err != nil {
			logger.Get().Errorw("Failed to get image information: "+err.Error(), append(w.loggerAttrs(), "err", err)...)
			return
		}

		logAttrs := append([]interface{}{"imageHash", imageHash}, w.loggerAttrs()...)
		if w.updateImageHash(imageHash) {
			logger.Get().Debugw("Image hash have been updated", logAttrs...)
			callback(targetName)
		}
	}
	go poll()

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
				poll()
			}
		}
	}()

	return w
}
