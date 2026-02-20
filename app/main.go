package main

import (
	"fmt"
	"os"
	"github.com/acskii/codecrafter-http-server/app/service"
)

func main() {
	if (len(os.Args) > 1 && os.Args[0] == "--directory") {
		fmt.Println("Detected flag --directory")
		fmt.Printf("Directory: %s\n", os.Args[2])
	}

	service.StartServer()
}