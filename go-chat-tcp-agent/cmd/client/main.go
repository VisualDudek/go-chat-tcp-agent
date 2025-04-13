package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chat-tcp-agent/internal/chat"
)

func main() {
	// Parse command line flags
	host := flag.String("host", "localhost", "Server host")
	port := flag.Int("port", 1234, "Server port")
	username := flag.String("username", "", "Your username")
	flag.Parse()

	// Validate username
	if *username == "" {
		fmt.Println("Error: Username is required. Use --username flag.")
		flag.Usage()
		os.Exit(1)
	}

	// Create client
	client := chat.NewClient(*username)

	// Connect to server
	serverAddr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Connecting to %s...\n", serverAddr)
	if err := client.Connect(serverAddr); err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// Create context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a goroutine
	go func() {
		<-signalChan
		fmt.Println("\nDisconnecting from server...")
		cancel()
	}()

	// Start client and handle errors
	if err := client.Start(ctx); err != nil {
		fmt.Printf("Client error: %v\n", err)
	}

	fmt.Println("Client disconnected")
}
