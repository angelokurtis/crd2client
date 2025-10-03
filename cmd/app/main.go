package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "Application exited with error", tint.Err(err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Application exited")
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})))

	undo, err := maxprocs.Set()
	defer undo()

	if err != nil {
		return fmt.Errorf("failed to set GOMAXPROCS: %w", err)
	}

	runner, cleanup, err := NewRunner(ctx)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer stop()
		return runner.Run(ctx)
	})

	g.Go(func() error {
		defer cleanup()

		<-ctx.Done()

		return nil
	})

	return g.Wait()
}

type Runner interface {
	Run(ctx context.Context) error
}

func NewRunner(ctx context.Context) (Runner, func(), error) {
	// TODO: configure the IoC container at this point and establish all necessary bindings
	return nil, func() {}, errors.New("not implemented")
}
