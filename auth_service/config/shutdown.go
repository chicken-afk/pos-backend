package config

import (
	"context"
	"log"

	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(init *Initialization, cancel context.CancelFunc, timeout time.Duration) {
	ctx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	log.Println("🔔 Received shutdown signal, shutting down...")

	// Cancel untuk worker/consumer jika ada
	if cancel != nil {
		cancel()
	}

	// Close DB, MQ, Elastic
	init.Close()

	// Shutdown Fiber jika ada
	if init.App != nil {
		if err := init.App.ShutdownWithContext(ctx); err != nil {
			log.Println("❌ Failed to shutdown Fiber cleanly:", err)
		} else {
			log.Println("✅ Server shutdown complete")
		}
	}
}
