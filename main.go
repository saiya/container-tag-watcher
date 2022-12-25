package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/saiya/container_tag_watcher/config"
	cr "github.com/saiya/container_tag_watcher/containerregistryclient"
	"github.com/saiya/container_tag_watcher/logger"
	"github.com/saiya/container_tag_watcher/notifier"
	"github.com/saiya/container_tag_watcher/watcher"
)

var crSetting = cr.Settings{}
var debugFlag = flag.Bool("debug", false, "shows DEBUG logs")

func main() {
	flag.BoolVar(&crSetting.EnableAwsEcrSupport, "aws-ecr", false, "enables AWS ECR credential handling")

	err := mainImpl()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func mainImpl() error {
	ctx := context.Background()

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		return fmt.Errorf("must give configuration file path in command line argument")
	}

	if *debugFlag {
		logger.EnableDebugLog()
	}

	cfgInput, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("cannot open configuration file \"%s\": %w", args[0], err)
	}
	defer cfgInput.Close()

	cfg, err := config.ParseConfig(cfgInput)
	if err != nil {
		return fmt.Errorf("invalid configuration file: %w", err)
	}

	n := notifier.NewNotifier(ctx, cfg)
	defer n.Close()

	w := watcher.NewWatcher(ctx, cfg, cr.Init(&crSetting), n.OnImageUpdated)
	defer w.Close()

	<-ctx.Done()
	return nil // Unreachable
}
