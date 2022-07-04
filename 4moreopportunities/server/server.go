package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/diyliv/grpc/4moreopportunities/ex/interceptors"
	"google.golang.org/grpc"
)

type server struct{}

func NewServer() *server {
	return &server{}
}

func (s *server) Run() {
	log.Printf("Starting gRPC server on port :50051")
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Error while listening: %v\n", err)
	}

	interceptors := interceptors.NewInterceptor()

	serv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryInterceptor))

	go func() {
		if err := serv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-ch

	serv.GracefulStop()
	log.Printf("Exiting was successful")
}
