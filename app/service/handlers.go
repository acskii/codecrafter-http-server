package service

import (
	"fmt"
	"net"
	"strings"
	"path/filepath"
	"os"
)

type Handler func(conn net.Conn, req *Request) bool

func RootHandler(conn net.Conn, req *Request) bool {
	if (req.Method == "GET" && req.Path == "/") {
		sendResponse(conn, req, "200 OK", "text/plain", "")
		return true
	}

	return false
}

func EchoHandler(conn net.Conn, req *Request) bool {
	if (req.Method == "GET" && strings.HasPrefix(req.Path, "/echo/")) {
		echoed := req.Path[len("/echo/"):]
		sendResponse(conn, req, "200 OK", "text/plain", echoed + "\n")
		return true
	}

	return false
}

func UserAgentHandler(conn net.Conn, req *Request) bool {
	if (req.Method == "GET" && strings.HasPrefix(req.Path, "/user-agent")) {
		agent := req.Headers["User-Agent"]
		sendResponse(conn, req, "200 OK", "text/plain", agent + "\n")
		return true
	}
	return false
}

func FileHandler(conn net.Conn, req *Request) bool {
	if (strings.HasPrefix(req.Path, "/files/")) {
		file := req.Path[len("/files/"):]

		if (getOSArg(1) == "--directory" && getOSArg(2) != "") {
			if (req.Method == "GET") {
				path := filepath.Join(getOSArg(2), file)
				content, err := os.ReadFile(path)
				if err != nil {
					fmt.Println("Failed to read file")
					sendResponse(conn, req, "404 Not Found", "text/plain", "")
				} else {
					sendResponse(conn, req, "200 OK", "application/octet-stream", string(content))
				}
			} else if (req.Method == "POST") {
				content := req.Body
				path := filepath.Join(getOSArg(2), file)
				data := []byte(content)
				err := os.WriteFile(path, data, 0644)

				if err != nil {
					fmt.Println("Failed to write file")
					sendResponse(conn, req, "403 Bad Request", "text/plain", "")
				} else {
					sendResponse(conn, req, "201 Created", "text/plain", "")
				}
			}
		} else {
			fmt.Println("No directory variable given. Aborting..")
		}
	}
	return false
}

func OtherwiseHandler(conn net.Conn, req *Request) bool {
	sendResponse(conn, req, "404 Not Found", "text/plain", "")
	return true
}