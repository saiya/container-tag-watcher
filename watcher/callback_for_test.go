package watcher_test

import (
	"time"

	"github.com/saiya/container_tag_watcher/watcher"
)

type callbackHistory struct {
	called chan interface{}
}

func newCallbackHistory() *callbackHistory {
	return &callbackHistory{
		called: make(chan interface{}, 1),
	}
}

func (h *callbackHistory) ImageUpdateCallback() watcher.ImageUpdateCallback {
	return func(targetName string) {
		h.called <- nil
	}
}

func (h *callbackHistory) Await(count int, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for i := 0; i < count; i++ {
		select {
		case <-timer.C:
			return false
		case <-h.called:
			break
		}
	}
	return true
}
