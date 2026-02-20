package main

import (
	"fmt"
	"os"
	"github.com/codecrafters-io/http-server-starter-go/app/service"
)

func main() {
	if (len(os.Args) > 1 && os.Args[0] == "--directory") {
		fmt.Println("Detected flag --directory")
		fmt.Printf("Directory: %s\n", os.Args[2])
	}

	service.StartServer()
}