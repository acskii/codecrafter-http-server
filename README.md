[![progress-banner](https://backend.codecrafters.io/progress/http-server/f19d9c30-fb04-4196-b6bb-fd8c95c8e1df)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)


# HTTP/1.1 Server

This is my implementation for the ["Build Your Own HTTP Server"](https://app.codecrafters.io/courses/http-server/overview) challenge, provided by [codecrafters.io](https://codecrafters.io), using Go's `net` package.

## Running the program

1. Ensure you have `go (1.25)` installed locally.
2. Run `./your_program.sh --directory /tmp/` to start the server (the `--directory` flag is required for file routes).

## Concepts Covered

### 1. Concurrent Request Handling
Unlike a sequential server that handles one client at a time, this implementation uses **Goroutines**. By spawning a new routine for every `ln.Accept()`, the server can process multiple TCP connections in parallel without blocking the main listener.


### 2. HTTP/1.1 Persistence (Keep-Alive)
I implemented a `for` loop within each connection handler to support **persistent connections**. This allows a single client (like `curl --next`) to send multiple requests over a single TCP socket.


### 3. Manual Protocol Parsing
Instead of using high-level frameworks, this project manually parses the raw byte stream from the socket:
* **Request Line:** Extracting Method, Path, and Protocol.
* **Headers:** Parsed into a `map[string]string`.
* **Body Extraction:** Utilizing `\r\n\r\n` as the delimiter between metadata and the message payload.


### 4. Connection Lifecycle Management
The server respects the `Connection: close` header. When detected, the server acknowledges the request by adding the header to the response and breaking the persistence loop.


## Features supported
- **HTTP Echo:** Returns the string provided in the URL path. Use route `/echo/{text}`.
- **User-Agent Detection:** Extracts and echoes the client's User-Agent header. Use route `/user-agent`.
- **File Server (GET):** Reads and returns content of a specified file from the directory variable set. Use route `/files/{filename}`.
- **File Upload (POST):** Receives data from the request body and writes it to file within the directory variable set. Use route `/files/{filename}`.
- **Concurrent Processing:** Handles multiple terminals simultaneously.
- **Persistent Connections:** Reuses TCP sockets for sequential requests.

## Testing with Curl

**Test Persistence:**
```bash
curl -v http://localhost:4221/echo/first --next http://localhost:4221/echo/second
```