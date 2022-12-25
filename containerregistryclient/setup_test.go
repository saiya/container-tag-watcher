package containerregistryclient_test

import (
	"testing"

	cr "github.com/saiya/container_tag_watcher/containerregistryclient"
	"github.com/saiya/container_tag_watcher/logger"
)

var client = cr.Init(&cr.Settings{})

func TestMain(m *testing.M) {
	logger.EnableDebugLog()
	m.Run()
}
