package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chat-tcp-agent/internal/chat"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 1234, "Port to listen on")
	flag.Parse()

	// Configure logging
	logger := log.New(os.Stdout, "[Server] ", log.LstdFlags)

	// Create server with address
	addr := fmt.Sprintf(":%d", *port)
	server := chat.NewServer(addr, logger)

	// Create context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.Start(ctx); err != nil {
			logger.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for interrupt signal
	<-signalChan
	logger.Println("Received interrupt signal, shutting down...")
	cancel()

	logger.Println("Server shutdown complete")
}
