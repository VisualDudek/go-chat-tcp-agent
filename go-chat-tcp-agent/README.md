# Go TCP Chat Application

A simple terminal-based TCP chat server and client implementation in Go.

## Features

- TCP-based client-server architecture
- Multiple concurrent client connections
- User join/leave notifications
- Real-time message broadcasting
- Clean shutdown handling
- Cross-platform compatibility

## Project Structure

```
go-chat-tcp-agent/
├── cmd/
│   ├── client/        # Client executable
│   └── server/        # Server executable
└── internal/
    └── chat/          # Core chat functionality
```

## Requirements

- Go 1.15 or later

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-chat-tcp-agent.git
   cd go-chat-tcp-agent
   ```

2. Initialize Go modules (if not already done):
   ```bash
   go mod init github.com/yourusername/go-chat-tcp-agent
   go mod tidy
   ```

### Building the Application

Build both the server and client at once from the project root:
```bash
go build -o server ./cmd/server
go build -o client ./cmd/client
```

Or build them separately:
```bash
# Build the server
cd cmd/server
go build -o ../../server

# Build the client
cd cmd/client
go build -o ../../client
```

## Usage Guide

### Starting the Server

Run the server with the default port (1234):
```bash
./server
```

Or specify a custom port:
```bash
./server --port 8080
```

The server will display a message showing the address it's listening on:
```
[Server] 2025/04/13 12:34:56 Server started on 0.0.0.0:8080
```

### Connecting with Clients

To connect a client to the server, open a new terminal window and run:

```bash
./client --username YourName
```

This connects to a server on localhost using the default port (1234).

For a server on a different host or port:
```bash
./client --host 192.168.1.100 --port 8080 --username Alice
```

### Chat Commands and Interaction

1. After connecting, you'll see a welcome message and join notification.
2. Type your message and press Enter to send it to all connected users.
3. Messages from other users will appear with their username: `[Bob]: Hello!`
4. System notifications (like users joining/leaving) appear with an asterisk: `* Charlie has joined the chat`
5. Press Ctrl+C to disconnect from the server.

### Running on Different Networks

To chat between computers on different networks:

1. The server must be running on a machine with a public IP or port forwarding enabled
2. Clients should connect using the public IP address of the server
3. Make sure the port is open in your firewall

## Common Issues and Troubleshooting

### Connection Refused

If you see "connection refused" errors when connecting:
- Verify the server is running
- Check that the host and port match the server configuration
- Ensure no firewall is blocking the connection

### Unexpected Disconnection

If clients disconnect unexpectedly:
- Check for network stability issues
- The server may have shut down
- Try reconnecting with the same command

### Multiple Users with Same Username

Currently, the application allows multiple users to use the same username. To avoid confusion, each user should pick a unique username.

## Development and Contributing

To make changes to the application:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

### Running Tests

Run tests from the project root:
```bash
go test ./...
```

## License

[MIT License](LICENSE)