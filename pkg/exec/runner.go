package exec

import (
	"context"
	"os"
	"os/exec"
	"sync/atomic"
	"time"
)

var duration = 200 * time.Millisecond

type Runner struct {
	command string
	args    []string
	timer   *time.Timer
	ctx     context.Context
	cancel  context.CancelFunc
	running atomic.Value
}

func NewRunner(command string, args ...string) *Runner {
	return &Runner{
		command: command,
		args:    args,
	}
}

func (r *Runner) Run(ctx context.Context) {
	if r.timer == nil {
		r.timer = time.NewTimer(duration)
	} else {
		if val, ok := r.running.Load().(bool); ok {
			if val {
				r.cancel()
				r.cancel = nil
			}
		}
		r.timer.Reset(duration)
	}

	r.ctx, r.cancel = context.WithCancel(ctx)
	go func() {
		for {
			r.running.Store(false)

			select {
			case <-r.timer.C:
				r.running.Store(true)
				cmd := exec.CommandContext(r.ctx, r.command, r.args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				if err := cmd.Run(); err != nil {
					continue
				}

			case <-ctx.Done():
				r.cancel()
				r.running.Store(false)
				return
			}
		}
	}()
}
