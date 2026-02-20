package service

import (
	"os"
	"net"
	"strings"
	"fmt"
)

type Request struct {
    Method  string
    Path    string
    Headers map[string]string
    Body    string
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

func getOSArg(index int) string {
	if index < 0 || index >= len(os.Args) {
		return ""
	}
	return os.Args[index]
}