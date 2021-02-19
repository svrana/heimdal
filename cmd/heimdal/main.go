package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattmoor/dep-notify/pkg/graph"
	"github.com/svrana/heimdal/pkg/exec"
)

func start(ctx context.Context, runners map[string]*exec.Runner) {
	fs := make(chan string)

	g, errCh, err := graph.New(func(ss graph.StringSet) {
		for target := range runners {
			if ss.Has(target) {
				fs <- target
			}
		}
	})
	if err != nil {
		fmt.Printf("failed to start: %s\n", err)
		return
	}
	defer g.Shutdown()

	for file := range runners {
		fmt.Printf("Watching %s and its deps for changes\n", file)
		if err := g.Add(file); err != nil {
			fmt.Printf("failed to start watch on: %s, %v\n", file, err)
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
			fmt.Printf("received: %s\n", file)
			if runner, ok := runners[file]; ok {
				runner.Run(ctx)
			}

		case err := <-errCh:
			fmt.Printf("got err: %v\n", err)
			return
		case <-sigint:
			return
		case <-ctx.Done():
			fmt.Printf("all done\n")
			return
		}
	}
}

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Printf("Must specify a path or package to watch\n")
		os.Exit(1)
	}

	watches := os.Args[1:]

	runners := make(map[string]*exec.Runner)
	for _, filename := range watches {
		runners[filename] = exec.NewRunner("make")
	}
	start(context.Background(), runners)
}
