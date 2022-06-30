package main

import (
	musicpb "github.com/diyliv/grpc/startWorkWithgRPC/ex/proto/music"
	"github.com/diyliv/grpc/startWorkWithgRPC/ex/server"
)

func main() {
	storage := make(map[string]musicpb.Music)
	server := server.NewServer(storage)
	server.Run()
}
