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

func routingGET(conn net.Conn, route string) {
	switch {
	case strings.HasPrefix(route, "/echo/"):
		handleEcho(conn, route[len("/echo/"):])
	case route == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default: 
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
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
			// Extract the request line
			data := string(buffer[:n])
			requestLine := strings.SplitN(data, "\r\n", 2)[0]
			parts := strings.Split(requestLine, " ")

			if (len(parts) > 0) {
				method := parts[0]
				// obtain request target
				target := parts[1]

				if (method == "GET") {
					// Ensure the method is GET
					// check for target and send response
					routingGET(c, target)
				}
			} else {
				c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
			c.Close()
		}(conn)
	}
}
