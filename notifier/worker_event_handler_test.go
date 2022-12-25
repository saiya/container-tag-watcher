package notifier_test

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/saiya/container_tag_watcher/notifier"
	"github.com/saiya/container_tag_watcher/testutil"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	history := newCommandsHistory()
	h := &notifier.WorkerEventHandler{
		TargetName: "test", EventName: "test",
		BacklogLimit: 1, ContinueOnError: false,
		Commands: [][]string{
			{"/bin/echo", "a", "b"},
			{"/bin/echo", "c", "d"},
		},
		CmdExecFn: history.CmdExecFn,
	}

	h.Start(testutil.Context(t))
	assert.True(t, h.Enqueue())

	assert.True(t, history.AwaitCmd(len(h.Commands), 3*time.Second))
	history.AssertEquals(t, h.Commands)
}

type TestContinueOnError_TestCase struct {
	ContinueOnError bool
	ExitCode        int
	Err             error
}

func TestContinueOnError(t *testing.T) {
	for _, c := range []struct {
		ExitCode        int
		Err             error
		ContinueOnError bool
	}{
		{-1, errors.New("test error"), false},
		{-1, errors.New("test error"), true},
		{128, &exec.ExitError{}, false},
		{128, &exec.ExitError{}, true},
	} {
		t.Run(
			fmt.Sprintf("ExitCode: %v, ContinueOnError: %v", c.ExitCode, c.ContinueOnError),
			func(t *testing.T) {
				ctx, ctxClose := context.WithCancel(testutil.Context(t))
				defer ctxClose()

				history := newCommandsHistory()

				var cmdCalled int32
				h := &notifier.WorkerEventHandler{
					TargetName: "test", EventName: "test",
					BacklogLimit: 1, ContinueOnError: c.ContinueOnError,
					Commands: [][]string{{"1"}, {"this-command-fails"}, {"3"}},
					CmdExecFn: func(name string, arg ...string) (int, error) {
						history.CmdExecFn(name, arg...)
						if atomic.AddInt32(&cmdCalled, 1) == 2 {
							return c.ExitCode, c.Err
						} else {
							return 0, nil
						}
					},
				}

				h.Start(ctx)
				assert.True(t, h.Enqueue())
				if c.ContinueOnError {
					assert.True(t, history.AwaitCmd(len(h.Commands), 3*time.Second))
					history.AssertEquals(t, h.Commands)
				} else {
					assert.True(t, history.AwaitCmd(len(h.Commands)-1, 3*time.Second))
					history.AssertEquals(t, h.Commands[0:2])
				}
			},
		)
	}
}

type commandsHistory struct {
	m       sync.Mutex
	history [][]string

	historyAppended chan interface{}
}

func newCommandsHistory() *commandsHistory {
	return &commandsHistory{
		history:         make([][]string, 0, 8),
		historyAppended: make(chan interface{}, 1),
	}
}

func (h *commandsHistory) CmdExecFn(name string, arg ...string) (int, error) {
	h.m.Lock()
	defer h.m.Unlock()

	h.history = append(h.history, append([]string{name}, arg...))
	h.historyAppended <- nil
	return 0, nil
}

func (h *commandsHistory) AwaitCmd(count int, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for i := 0; i < count; i++ {
		select {
		case <-timer.C:
			return false
		case <-h.historyAppended:
			break
		}
	}
	return true
}

func (h *commandsHistory) AssertEquals(t *testing.T, expected [][]string) {
	h.m.Lock()
	defer h.m.Unlock()

	if !assert.Equal(t, len(expected), len(h.history)) {
		return
	}
	for i, expectedCmd := range expected {
		assert.Equal(t, expectedCmd, h.history[i])
	}
}
