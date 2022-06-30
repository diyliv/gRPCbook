package main

import (
	"context"
	"log"
	"time"

	musicpb "github.com/diyliv/grpc/startWorkWithgRPC/ex/proto/music"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("not available: %v\n", err)
	}

	defer cc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := musicpb.NewMusicServiceClient(cc)
	resp, err := c.AddMusic(ctx, &musicpb.MusicRequest{
		&musicpb.Music{
			SongName: "The Amity Affliction - IVY",
			AuthorInfo: &musicpb.Author{
				Name:      "Ahren Stringer, Joel Birch, Dan Brown, Joe Longobardi",
				Age:       "NULL",
				Followers: "NULL",
				Tracks:    "enough",
			},
			Description:  "cool song",
			SongDuration: timestamppb.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Error while calling Add RPC: %v\n", err)
	}

	getInfo, err := c.GetMusic(ctx, &musicpb.MusicID{Id: resp.GetId()})
	if err != nil {
		log.Fatalf("Error while calling Get RPC: %v\n", err)
	}

	log.Printf("%v\n", getInfo)
}
