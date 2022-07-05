package server

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/diyliv/grpc/startWorkWithgRPC/ex/controllers"
	musicpb "github.com/diyliv/grpc/startWorkWithgRPC/ex/proto/music"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	musicInfo map[string]musicpb.Music // на данный момент, я не хочу подключать базу данных и буду хранить данные в виде ключ:значение
}

func NewServer(musicInfo map[string]musicpb.Music) *server {
	return &server{musicInfo: musicInfo}
}

var (
	pemFile = "../certs/leaf.pem"
	keyFile = "../certs/leaf.key"
)

func (s *server) Run() {
	log.Printf("Starting gRPC server on port :50051")
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Error while starting gRPC server: %v\n", err)
	}

	cert, err := tls.LoadX509KeyPair(pemFile, keyFile)
	if err != nil {
		log.Fatalf("Error while parsing public/private key: %v\n", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	serv := grpc.NewServer(opts...)
	controllers := controllers.NewController(s.musicInfo)
	musicpb.RegisterMusicServiceServer(serv, controllers)

	go func() {
		if err := serv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	serv.GracefulStop()
	log.Printf("Exiting was successful")
}
