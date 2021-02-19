package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattmoor/dep-notify/pkg/graph"
)

func start(ctx context.Context, files []string) {
	fs := make(chan string)

	g, errCh, err := graph.New(func(ss graph.StringSet) {
		// ugh, a write in nvim ends up with 4 events for the same
		// file, remove, create, chmod write and chmod so will have to
		// buffer these
		fmt.Printf("event: %v\n ", ss)
		for _, file := range files {
			if ss.Has(file) {
				fs <- file
			}
		}
	})
	if err != nil {
		fmt.Printf("failed to start: %s\n", err)
		return
	}
	defer g.Shutdown()

	for _, file := range files {
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
	start(context.Background(), os.Args[1:])
}
