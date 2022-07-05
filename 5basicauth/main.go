package main

import (
	"github.com/diyliv/grpc/5security/5auth/server"
)

func main() {
	server := server.NewServer()
	server.Run()
}
