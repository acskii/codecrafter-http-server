package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"path/filepath"
)

type Request struct {
    Method  string
    Path    string
    Headers map[string]string
    Body    string
}

func handleEcho(conn net.Conn, req *Request, s string) {
	sendResponse(conn, req, "200 OK", "text/plain", s + "\n")
}

func handleUserAgent(conn net.Conn, req *Request) {
	agent := req.Headers["User-Agent"]
	sendResponse(conn, req, "200 OK", "text/plain", agent + "\n")
}

func handleFiles(conn net.Conn, req *Request, file string) {
	if os.Args[1] == "--directory" {
		path := filepath.Join(os.Args[2], file)
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Failed to read file")
			sendResponse(conn, req, "404 Not Found", "text/plain", "")
		} else {
			sendResponse(conn, req, "200 OK", "application/octet-stream", string(content))
		}
	}
}

func handlePostFiles(conn net.Conn, req *Request, file string) {
	if os.Args[1] == "--directory" {
		content := req.Body
		path := filepath.Join(os.Args[2], file)
		data := []byte(content)
		err := os.WriteFile(path, data, 0644)

		if err != nil {
			fmt.Println("Failed to write file")
			sendResponse(conn, req, "403 Bad Request", "text/plain", "")
		} else {
			sendResponse(conn, req, "201 Created", "text/plain", "")
		}
	}
}

func routingGET(conn net.Conn, req *Request) {
	route := req.Path

	switch {
	case strings.HasPrefix(route, "/echo/"):
		handleEcho(conn, req, route[len("/echo/"):])
	case strings.HasPrefix(route, "/user-agent"):
		handleUserAgent(conn, req)
	case strings.HasPrefix(route, "/files/"):
		handleFiles(conn, req, route[len("/files/"):])
	case route == "/":
		sendResponse(conn, req, "200 OK", "text/plain", "")
	default: 
		sendResponse(conn, req, "404 Not Found", "text/plain", "")
	}
}

func routingPOST(conn net.Conn, req *Request) {
	route := req.Path

	switch {
	case strings.HasPrefix(route, "/files/"):
		handlePostFiles(conn, req, route[len("/files/"):])
	default: 
		conn.Write([]byte("HTTP/1.1 403 Bad Request\r\n\r\n"))
	}
}

func parseRequest(request string) *Request {
	parts := strings.Split(request, "\r\n")

	if (len(parts) == 0) {
		return nil
	}

	// Extract request line
	requestLine := strings.Split(parts[0], " ")
	method := requestLine[0]
	path := requestLine[1]

	// Extract request headers
	headers := make(map[string]string)

	// Skip the request line
	for i := 1; i < len(parts); i++ {
		line := parts[i]

		// Stop at \r\n\r\n
		if line == "" {
			break
		}

		// Get header and value
		l := strings.SplitN(line, ": ", 2)
		if len(l) == 2 {
			headers[l[0]] = l[1]
		}
	}

	// Extract request body
	body := strings.SplitN(request, "\r\n\r\n", 2)[1]

	// Return Request object
	return &Request {
		Method:  method,
		Path:    path,
		Headers: headers,
		Body:	 body,
	}
}

func sendResponse(conn net.Conn, req *Request, status string, contentType string, body string) {
    response := fmt.Sprintf("HTTP/1.1 %s\r\n", status)
    response += fmt.Sprintf("Content-Type: %s\r\n", contentType)
    response += fmt.Sprintf("Content-Length: %d\r\n", len(body))

	if (req.Headers["Connection"] == "close") {
		response += "Connection: close\r\n"
	}

	response += "\r\n"
    response += body
    conn.Write([]byte(response))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	
	for {
		// Read the request
		buffer := make([]byte, 1024)
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
			if (req.Method == "GET") {
				routingGET(conn, req)
			} else if (req.Method == "POST") {
				routingPOST(conn, req)
			}

			if (req.Headers["Connection"] == "close") {
				break
			}
		} else {
			sendResponse(conn, req, "404 Not Found", "text/plain", "")
			return
		}
	}
}

func main() {
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
