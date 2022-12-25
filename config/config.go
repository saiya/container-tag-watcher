package config

import (
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets map[string]*WatchTarget `json:"targets" yaml:"targets"`
}

func (cfg *Config) TargetNames() []string {
	keys := make([]string, 0, len(cfg.Targets))
	for k := range cfg.Targets {
		keys = append(keys, k)
	}
	return keys
}

type WatchTarget struct {
	Platform        string `json:"platform" yaml:"platform"`
	PollingInterval string `json:"polling-interval" yaml:"polling-interval"`
	ContinueOnError bool   `json:"continue-on-error" yaml:"continue-on-error"`
	BacklogLimit    *int   `json:"backlog-limit" yaml:"backlog-limit"`

	Commands []interface{} `json:"commands" yaml:"commands"` // After post-process, elements are always []string
}

const PlatformDefault = "linux/amd64"
const PollingIntervalDefault = "3m"
const ContinueOnErrorDefault = false
const NotificationBacklogLimitDefault = 1

func ParseConfig(in io.Reader) (*Config, error) {
	src, err := io.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := Config{}
	err = yaml.Unmarshal(src, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file YAML: %w", err)
	}

	var parseError *ConfigParseError
	for key, target := range config.Targets {
		parseError = mergeParseError(parseError, target.postProcess(fmt.Sprintf("targets.\"%s\"", key)))
	}
	if parseError == nil { // https://go.dev/doc/faq#nil_error
		return &config, nil
	} else {
		return &config, parseError
	}
}

func (target *WatchTarget) postProcess(location string) *ConfigParseError {
	var parseError *ConfigParseError

	if target.Platform == "" {
		target.Platform = PlatformDefault
	}

	if target.PollingInterval == "" {
		target.PollingInterval = PollingIntervalDefault
	}
	if _, err := time.ParseDuration(target.PollingInterval); err != nil {
		parseError = appendParseError(parseError, &ConfigParseErrorDetail{
			Location: location + ".polling-interval",
			Message:  fmt.Sprintf("cannot parse given duration string: %v", err),
		})
	}

	if target.BacklogLimit == nil {
		limit := NotificationBacklogLimitDefault
		target.BacklogLimit = &limit
	}

	// Normalize elements of target.Commands to []string
	for i, command := range target.Commands {
		if str, ok := command.(string); ok {
			target.Commands[i] = []string{"sh", "-c", str}
		} else if arr, ok := command.([]interface{}); ok {
			strArr := make([]string, len(arr))
			target.Commands[i] = strArr
			for j, arg := range arr {
				if strArg, ok2 := arg.(string); ok2 {
					strArr[j] = strArg
				} else {
					parseError = appendParseError(parseError, &ConfigParseErrorDetail{
						Location: fmt.Sprintf("%s.commands[%d][%d]", location, i, j),
						Message:  "must be string (not list, map, ...)",
					})
				}
			}
		} else {
			parseError = appendParseError(parseError, &ConfigParseErrorDetail{
				Location: fmt.Sprintf("%s.commands[%d]", location, i),
				Message:  "must be string or array of string",
			})
		}
	}

	return parseError
}
