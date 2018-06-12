package main

import (
	"github.com/gambarini/grpcdemo/svcutils"
	"github.com/gambarini/grpcdemo/testsvc/internal/server"
)

func main() {

	main := svcutils.Main{
		ServerPort: 50051,
		Name:       "Test Service",
		Server:     &server.TestServer{},
	}

	main.Run()
}