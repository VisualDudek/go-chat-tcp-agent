package chat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// ServerClient represents a client connection in the server
type ServerClient struct {
	ID       string
	Username string
	Conn     net.Conn
	Writer   *bufio.Writer
	Server   *Server
}

// Server represents the chat server
type Server struct {
	Addr       string
	Clients    map[string]*ServerClient
	ClientsMu  sync.RWMutex
	Logger     *log.Logger
	ListenAddr string
}

// NewServer creates a new chat server
func NewServer(addr string, logger *log.Logger) *Server {
	if logger == nil {
		logger = log.Default()
	}

	return &Server{
		Addr:    addr,
		Clients: make(map[string]*ServerClient),
		Logger:  logger,
	}
}

// Start starts the chat server
func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer listener.Close()

	s.ListenAddr = listener.Addr().String()
	s.Logger.Printf("Server started on %s", s.ListenAddr)

	go func() {
		<-ctx.Done()
		s.Logger.Println("Shutting down server...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				s.Logger.Printf("Error accepting connection: %v", err)
				continue
			}
		}

		go s.handleConnection(conn)
	}
}

// handleConnection processes a new client connection
func (s *Server) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// First line should be the username
	username, err := reader.ReadString('\n')
	if err != nil {
		s.Logger.Printf("Error reading username: %v", err)
		conn.Close()
		return
	}

	// Trim newline character
	username = username[:len(username)-1]

	// Create new client
	client := &ServerClient{
		ID:       conn.RemoteAddr().String(),
		Username: username,
		Conn:     conn,
		Writer:   writer,
		Server:   s,
	}

	// Add client to the server
	s.addClient(client)

	// Broadcast join message
	s.Broadcast(fmt.Sprintf("* %s has joined the chat", client.Username), client.ID)

	// Start handling client messages
	s.handleClient(client, reader)
}

// addClient adds a client to the server
func (s *Server) addClient(client *ServerClient) {
	s.ClientsMu.Lock()
	defer s.ClientsMu.Unlock()

	s.Clients[client.ID] = client
	s.Logger.Printf("Client connected: %s (%s)", client.Username, client.ID)
}

// removeClient removes a client from the server
func (s *Server) removeClient(clientID string) {
	s.ClientsMu.Lock()
	defer s.ClientsMu.Unlock()

	if client, exists := s.Clients[clientID]; exists {
		delete(s.Clients, clientID)
		client.Conn.Close()
		s.Logger.Printf("Client disconnected: %s (%s)", client.Username, client.ID)
	}
}

// handleClient processes messages from a client
func (s *Server) handleClient(client *ServerClient, reader *bufio.Reader) {
	defer func() {
		s.removeClient(client.ID)
		s.Broadcast(fmt.Sprintf("* %s has left the chat", client.Username), client.ID)
	}()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				s.Logger.Printf("Error reading from client %s: %v", client.ID, err)
			}
			break
		}

		// Trim newline character
		message = message[:len(message)-1]

		// Format the message with the username
		formattedMsg := fmt.Sprintf("[%s]: %s", client.Username, message)

		// Broadcast the message to all clients
		s.Broadcast(formattedMsg, client.ID)

		// Log the message
		s.Logger.Printf("Message from %s: %s", client.Username, message)
	}
}

// Broadcast sends a message to all connected clients except the sender
func (s *Server) Broadcast(message string, senderID string) {
	s.ClientsMu.RLock()
	defer s.ClientsMu.RUnlock()

	for id, client := range s.Clients {
		// Don't send message back to sender if senderID is provided
		if senderID != "" && id == senderID {
			continue
		}

		go func(c *ServerClient, msg string) {
			if _, err := c.Writer.WriteString(msg + "\n"); err != nil {
				s.Logger.Printf("Error sending message to %s: %v", c.ID, err)
				return
			}
			if err := c.Writer.Flush(); err != nil {
				s.Logger.Printf("Error flushing message to %s: %v", c.ID, err)
			}
		}(client, message)
	}
}
