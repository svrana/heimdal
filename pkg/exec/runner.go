package exec

import (
	"context"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/svrana/heimdal/pkg/zapped"
)

var defaultDuration = 450 * time.Millisecond

type Runner struct {
	name    string
	command string
	args    []string
	timer   *time.Timer
	ctx     context.Context
	cancel  context.CancelFunc
	running atomic.Value
	delay   time.Duration
}

func NewRunner(name, command string, args ...string) *Runner {
	return &Runner{
		name:    name,
		command: command,
		args:    args,
		delay:   defaultDuration,
	}

}

func (r *Runner) WithDelayDuration(duration int) *Runner {
	r.delay = time.Duration(duration) * time.Millisecond
	return r
}

func (r *Runner) Run(ctx context.Context) {
	l := zapped.FromContext(ctx)

	if r.timer == nil {
		r.timer = time.NewTimer(r.delay)
	} else {
		if val, ok := r.running.Load().(bool); ok {
			if val {
				l.Debugf("Incoming event: canceling %s", r.name)
				r.cancel()
				r.cancel = nil
			}
		}
		l.Debugf("Delaying %s for %s", r.name, r.delay)
		r.timer.Reset(r.delay)
	}

	r.ctx, r.cancel = context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-r.timer.C:
				r.running.Store(true)
				cmd := exec.CommandContext(r.ctx, r.command, r.args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				l.Debugf("running command")
				if err := cmd.Run(); err != nil {
					if r.ctx.Err() != nil {
						l.Debugf("%s cancelled", r.name)
					}
				}
				r.running.Store(false)
			case <-ctx.Done():
				r.cancel()
				r.running.Store(false)
				return
			}
		}
	}()
}
