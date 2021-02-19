package heimdal

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattmoor/dep-notify/pkg/graph"

	"github.com/svrana/heimdal/pkg/exec"
	"github.com/svrana/heimdal/pkg/zapped"
)

func Start(ctx context.Context, runners map[string]*exec.Runner) {
	l := zapped.FromContext(ctx)
	fs := make(chan string)

	g, errCh, err := graph.New(func(ss graph.StringSet) {
		for target := range runners {
			if ss.Has(target) {
				fs <- target
			}
		}
	})
	if err != nil {
		l.Errorf("failed to start: %s", err)
		return
	}
	defer g.Shutdown()

	for file := range runners {
		l.Infof("Watching %s and its deps for changes", file)
		if err := g.Add(file); err != nil {
			l.Errorf("failed to start watch on: %s, %v", file, err)
			return
		}
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case file, ok := <-fs:
			if !ok {
				fs = nil
				break
			}
			l.Debugf("received: %s", file)
			if runner, ok := runners[file]; ok {
				runner.Run(ctx)
			}

		case err := <-errCh:
			l.Infof("got err: %v", err)
			return
		case <-sigint:
			return
		case <-ctx.Done():
			l.Infof("all done")
			return
		}
	}
}
