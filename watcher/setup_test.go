package watcher_test

import (
	"testing"

	"github.com/saiya/container_tag_watcher/logger"
)

func TestMain(m *testing.M) {
	logger.EnableDebugLog()
	m.Run()
}
