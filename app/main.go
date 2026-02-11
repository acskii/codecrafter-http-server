package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleEcho(conn net.Conn, s string) {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(s), s)
	conn.Write([]byte(response))
}

func handleUserAgent(conn net.Conn, headers map[string]string) {
	agent := headers["User-Agent"]
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(agent), agent)
	conn.Write([]byte(response))
}

func routingGET(conn net.Conn, route string, headers map[string]string) {
	switch {
	case strings.HasPrefix(route, "/echo/"):
		handleEcho(conn, route[len("/echo/"):])
	case strings.HasPrefix(route, "/user-agent"):
		handleUserAgent(conn, headers)
	case route == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default: 
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func getHeaders(request string) map[string]string {
	// Separate by \r\n
	lines := strings.Split(request, "\r\n")
	headers := make(map[string]string)

	// Skip the request line
	for i := 1; i < len(lines); i++ {
		line := lines[i]

		// Stop at \r\n\r\n
		if line == "" {
			break
		}

		// Get header and value
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}

	return headers
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", ":4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func(c net.Conn) {
			// Read the request
			buffer := make([]byte, 1024)
			n, err := c.Read(buffer)
			if (err != nil) {
				fmt.Println("Error reading request: ", err.Error())
				c.Close()
			}
			data := string(buffer[:n])
			// Extract the request line
			requestLine := strings.SplitN(data, "\r\n", 2)[0]
			// Extract the headers
			headers := getHeaders(data)

			parts := strings.Split(requestLine, " ")

			if (len(parts) > 0) {
				method := parts[0]
				// obtain request target
				target := parts[1]

				if (method == "GET") {
					// Ensure the method is GET
					// check for target and send response
					routingGET(c, target, headers)
				}
			} else {
				c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
			c.Close()
		}(conn)
	}
}
