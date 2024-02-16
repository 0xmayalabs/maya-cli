package main

import (
	"context"
	"github.com/twindevs/project-maya/cmd"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err := cmd.New().ExecuteContext(ctx)

	cancel()

	if err != nil {
		slog.Error("Fatal error", "err", err)
		os.Exit(1)
	}
}
