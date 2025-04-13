# Project Requirements Document

**Project Title:** Terminal TCP Chat Server  
**Author:** [Your Name]  
**Date:** 2025-04-13  
**Version:** 1.0

---

## 1. Purpose

The purpose of this project is to design and implement a **terminal-based TCP chat server** that allows multiple clients to connect, send, and receive messages in real time via the command line interface (CLI). The system should be simple, robust, and easily extendable for future features like authentication or message history.

---

## 2. Scope

- Server and client will communicate over TCP.
- Clients will use a terminal interface to send and receive messages.
- The server will handle multiple clients concurrently.
- Basic features like user join/leave notifications and message broadcasting will be supported.

---

## 3. Functional Requirements

### 3.1 Server Requirements

- Accept incoming TCP connections on a configurable port.
- Support multiple clients concurrently.
- Broadcast received messages to all connected clients.
- Notify clients when a new user joins or leaves.
- Handle clean shutdown (e.g., on SIGINT).
- Optionally: log messages to the console or a file.

### 3.2 Client Requirements

- Connect to the server via TCP.
- Read user input from the terminal and send it to the server.
- Receive and display broadcasted messages.
- Handle disconnects and server shutdown gracefully.

---

## 4. Non-Functional Requirements

- Written in a high-level language like Python or Go.
- Code should be modular and documented.
- System should handle at least 10 concurrent clients.
- Should function cross-platform (Linux, Windows, macOS).
- No GUI; purely terminal-based.

---

## 5. User Interface (CLI)

- **Start server:** `./server [--port 1234]`
- **Start client:** `./client --host 127.0.0.1 --port 1234 --username Alice`
- **Chat format:** `[Alice]: Hello world`
- **System notifications:** `* Alice has joined the chat`

---

## 6. Assumptions

- Users are running the application in a terminal.
- There’s no authentication or encryption (for MVP).
- Users have basic networking permissions and access.

---

## 7. Constraints

- Only standard libraries should be used for MVP.
- Server should not crash due to client misbehavior (e.g., abrupt disconnects).

---

## 8. Future Enhancements (Optional)

- Add username authentication.
- Implement private messaging.
- Add message history.
- Encrypt communication using TLS.
- WebSocket or GUI interface.

---

## 9. Deliverables

- Source code for server and client.
- README with setup and usage instructions.
- Test cases or manual testing plan.
- PRD and optional design diagram.
