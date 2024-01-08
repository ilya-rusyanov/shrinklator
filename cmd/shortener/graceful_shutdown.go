package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type shutdowner interface {
	Stop(context.Context) error
}

func gracefulShutdown(ctx context.Context, log Logger, targets ...shutdowner) <-chan struct{} {
	doneCh := make(chan struct{})

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-term
		for _, target := range targets {
			if err := target.Stop(ctx); err != nil {
				log.Fatalf("failed to stop: %w", err)
			}
		}
		close(doneCh)
	}()

	return doneCh
}
