package notifier

import (
	"context"
	"os/exec"
	"time"

	"github.com/saiya/container_tag_watcher/logger"
)

type WorkerEventHandler struct {
	TargetName string
	EventName  string

	BacklogLimit    int
	ContinueOnError bool
	Commands        [][]string

	// A function executes command and returns exitCode
	// (For testing) Can inject os/exec.Command.exec mock
	CmdExecFn func(name string, arg ...string) (int, error)

	ch chan interface{}
}

func (h *WorkerEventHandler) logAttrs() []interface{} {
	return []interface{}{
		"target", h.TargetName,
		"event", h.EventName,
	}
}

func (h *WorkerEventHandler) Enqueue() bool {
	select {
	case h.ch <- nil:
		logger.Get().Debugw(
			"Command execution queued",
			h.logAttrs()...,
		)
		return true
	default:
		logger.Get().Infow(
			"Command queue is full (= command still running for past events), ignored event",
			h.logAttrs()...,
		)
		return false
	}
}

func (h *WorkerEventHandler) Start(ctx context.Context) {
	h.ch = make(chan interface{}, h.BacklogLimit)
	if h.CmdExecFn == nil {
		h.CmdExecFn = func(name string, arg ...string) (int, error) {
			cmd := exec.Command(name, arg...)
			err := cmd.Run()
			return cmd.ProcessState.ExitCode(), err
		}
	}

	go func() {
	Loop:
		for {
			select {
			case <-ctx.Done():
				logger.Get().Debugw("Notification worker shutdown...", h.logAttrs()...)
				break Loop
			case <-h.ch:
				h.execCommands(ctx)
			}
		}
		logger.Get().Debugw("Notification worker closed", h.logAttrs()...)
	}()
}

func (h *WorkerEventHandler) execCommands(ctx context.Context) {
CmdLoop:
	for i, cmdStrs := range h.Commands {
		zapFields := append([]interface{}{
			"command", cmdStrs[0], // Avoid exposing arguments into log for safety
			"sequence", i,
		}, h.logAttrs()...)
		logger.Get().Infow("Executing command", zapFields...)

		startAt := time.Now()
		exitCode, err := h.CmdExecFn(cmdStrs[0], cmdStrs[1:]...)
		endAt := time.Now()
		zapFields = append(
			zapFields,
			"exitCode", exitCode,
			"startAt", startAt, "endAt", endAt, "seconds", endAt.Sub(startAt).Seconds(),
		)

		if err == nil {
			logger.Get().Infow("Command successfully ended", zapFields...)
		} else {
			if _, ok := err.(*exec.ExitError); ok {
				logger.Get().Warnw("Command ended in non-zero code", zapFields...)
			} else {
				zapFields = append(zapFields, "err", err)
				logger.Get().Errorw("Failed to execute command: "+err.Error(), zapFields...)
			}
			if !h.ContinueOnError {
				break CmdLoop
			}
		}
	}
}
