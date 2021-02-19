package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/svrana/heimdal/pkg/zapped"
)

var duration = 800 * time.Millisecond

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
	l := zapped.FromContext(ctx)

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
				fmt.Println("running command")
				if err := cmd.Run(); err != nil {
					if r.ctx.Err() == nil {
						l.Info("command cancelled")
					}
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
