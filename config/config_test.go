package config_test

import (
	"strings"
	"testing"

	"github.com/saiya/container_tag_watcher/config"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cfg, err := config.ParseConfig(strings.NewReader(`
targets:
  "saiya/container-tag-watcher:latest":
    commands:
      - "echo a b c"
  "12345678890.dkr.ecr.ap-northeast-1.amazonaws.com/my-application:latest":
    polling-interval: "123m"
    continue-on-error: true
    backlog-limit: 123
    commands:
      - "echo command #1"
      -
        - "/bin/echo"
        - "a b c"
  `))
	if err != nil {
		t.Fatal(err)
	}

	keys := make([]string, 0, len(cfg.Targets))
	for k := range cfg.Targets {
		keys = append(keys, k)
	}

	var target1 = "saiya/container-tag-watcher:latest"
	assert.Equal(t, config.PollingIntervalDefault, cfg.Targets[target1].PollingInterval)
	assert.Equal(t, config.ContinueOnErrorDefault, cfg.Targets[target1].ContinueOnError)
	assert.Equal(t, config.NotificationBacklogLimitDefault, *cfg.Targets[target1].BacklogLimit)
	assert.Equal(t, []string{"sh", "-c", "echo a b c"}, cfg.Targets[target1].Commands[0])

	var target2 = "12345678890.dkr.ecr.ap-northeast-1.amazonaws.com/my-application:latest"
	assert.Equal(t, "123m", cfg.Targets[target2].PollingInterval)
	assert.Equal(t, !config.ContinueOnErrorDefault, cfg.Targets[target2].ContinueOnError)
	assert.Equal(t, 123, *cfg.Targets[target2].BacklogLimit)
	assert.Equal(t, []string{"sh", "-c", "echo command #1"}, cfg.Targets[target2].Commands[0])
	assert.Equal(t, []string{"/bin/echo", "a b c"}, cfg.Targets[target2].Commands[1])
}

func TestInvalidCommand(t *testing.T) {
	_, err := config.ParseConfig(strings.NewReader(`
targets:
  "test/test:latest":
    commands:
      - invalid: map
      -
        - "/bin/sh"
        - invalid: map
  `))

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `targets."test/test:latest".commands[0]: must be string or array of string`)
	assert.Contains(t, err.Error(), `targets."test/test:latest".commands[1][1]: must be string`)
}

func TestInvalidPollingInterval(t *testing.T) {
	_, err := config.ParseConfig(strings.NewReader(`
  targets:
    "test/test:latest":
      polling-interval: "invalid"
      commands:
        - echo foo bar
    `))

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `targets."test/test:latest".polling-interval: cannot parse given duration string`)
}
