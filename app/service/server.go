package service

import (
	"fmt"
	"net"
	"os"
)

func StartServer() {
	fmt.Println("Server started!")

	l, err := net.Listen("tcp", ":4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
        os.Exit(1)
	}
	
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	
	for {
		// Read the request
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if (err != nil) {
			if (err.Error() != "EOF") {
				fmt.Println("Error reading request: ", err.Error())
			}
			return
		}
		data := string(buffer[:n])

		// Extract the request components
		req := parseRequest(data)

		if (req != nil) {
			delegantHandlers(conn, req)
			if (req.Headers["Connection"] == "close") {
				break
			}
		} else {
			sendResponse(conn, req, "404 Not Found", "text/plain", "")
			return
		}
	}
}

func delegantHandlers(conn net.Conn, req *Request) {
	h := []Handler{
		RootHandler,
		EchoHandler,
		UserAgentHandler,
		FileHandler,
		OtherwiseHandler,
	}

	for _, handler := range h {
		if handler(conn, req) {
			return
		}
	}
}