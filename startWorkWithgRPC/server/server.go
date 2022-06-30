package server

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/diyliv/grpc/startWorkWithgRPC/ex/controllers"
	musicpb "github.com/diyliv/grpc/startWorkWithgRPC/ex/proto/music"
	"google.golang.org/grpc"
)

type server struct {
	musicInfo map[string]musicpb.Music // на данный момент, я не хочу подключать базу данных и буду хранить данные в виде ключ:значение
}

func NewServer(musicInfo map[string]musicpb.Music) *server {
	return &server{musicInfo: musicInfo}
}
func (s *server) Run() {
	log.Printf("Starting gRPC server on port :50051")
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Error while starting gRPC server: %v\n", err)
	}

	serv := grpc.NewServer()
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
