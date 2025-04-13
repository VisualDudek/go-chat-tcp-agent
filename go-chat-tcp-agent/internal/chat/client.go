package chat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Client represents a chat client
type Client struct {
	Username string
	Conn     net.Conn
	Reader   *bufio.Reader
	Writer   *bufio.Writer
}

// NewClient creates a new chat client
func NewClient(username string) *Client {
	return &Client{
		Username: username,
		Reader:   bufio.NewReader(os.Stdin),
	}
}

// Connect connects the client to a chat server
func (c *Client) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	c.Conn = conn
	c.Writer = bufio.NewWriter(conn)

	// Send username as the first message
	if err := c.sendMessage(c.Username); err != nil {
		return fmt.Errorf("failed to send username: %w", err)
	}

	return nil
}

// Start starts the client's message loop
func (c *Client) Start(ctx context.Context) error {
	// Create a context that will be canceled when the client disconnects
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start a goroutine to receive messages from the server
	go c.receiveMessages(ctx, cancel)

	// Start the message input loop
	fmt.Println("Connected to chat server. Type your messages and press Enter to send.")
	fmt.Println("Press Ctrl+C to exit.")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Read user input
			fmt.Print("> ")
			input, err := c.Reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// User pressed Ctrl+D
					return nil
				}
				return fmt.Errorf("failed to read input: %w", err)
			}

			// Trim whitespace
			input = strings.TrimSpace(input)

			// Skip empty messages
			if input == "" {
				continue
			}

			// Send the message to the server
			if err := c.sendMessage(input); err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}
	}
}

// sendMessage sends a message to the server
func (c *Client) sendMessage(message string) error {
	_, err := c.Writer.WriteString(message + "\n")
	if err != nil {
		return err
	}

	return c.Writer.Flush()
}

// receiveMessages receives and displays messages from the server
func (c *Client) receiveMessages(ctx context.Context, cancel context.CancelFunc) {
	reader := bufio.NewReader(c.Conn)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Error reading from server: %v\n", err)
				}
				fmt.Println("\nDisconnected from server.")
				cancel() // Cancel the context to signal the main loop to exit
				return
			}

			// Print the message (remove the trailing newline)
			fmt.Print("\r" + strings.TrimSuffix(message, "\n") + "\n> ")
		}
	}
}

// Close closes the client connection
func (c *Client) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}
