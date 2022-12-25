package watcher_test

import (
	"testing"
	"time"

	"github.com/saiya/container_tag_watcher/config"
	"github.com/saiya/container_tag_watcher/testutil"
	"github.com/saiya/container_tag_watcher/watcher"
	"github.com/stretchr/testify/assert"
)

func TestWatcher(t *testing.T) {
	crClient := newStubClient()
	callbackHistory := newCallbackHistory()

	imageArch := "amd64/linux"
	imageName := "test/test:latest"
	crClient.SetImageDigest(imageArch, imageName, "1111")
	w := watcher.NewWatcher(
		testutil.Context(t),
		&config.Config{
			Targets: map[string]*config.WatchTarget{
				imageName: {
					Platform:        imageArch,
					PollingInterval: "10ms",
				},
			},
		},
		crClient,
		callbackHistory.ImageUpdateCallback(),
	)
	defer w.Close()

	time.Sleep(15 * time.Millisecond)
	crClient.SetImageDigest(imageArch, imageName, "2222")

	assert.True(t, callbackHistory.Await(1, time.Second))
}
