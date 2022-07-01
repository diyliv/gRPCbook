package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	musicpb "github.com/diyliv/grpc/startWorkWithgRPC/ex/proto/music"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	// serverStreaming(ctx, c)
	clientStreaming(ctx, c)
}

func serverStreaming(ctx context.Context, c musicpb.MusicServiceClient) {
	res, err := c.SearchMusic(ctx, &wrapperspb.StringValue{Value: "lofi/rain/forest"})
	if err != nil {
		log.Fatalf("Error while calling SearchMusic RPC: %v\n", err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving data: %v\n", err)
		}
		log.Printf("%v\n", msg.GetSearch())
	}

}

func clientStreaming(ctx context.Context, c musicpb.MusicServiceClient) {
	req := []musicpb.Music{
		musicpb.Music{
			SongName: "Bring Me The Horizon - sleepwalking",
			AuthorInfo: &musicpb.Author{
				Name:      "Oliver Sykes, Matt Kean, Lee Malia, Matt Nicholls, Jordan Fish",
				Age:       "NULL",
				Followers: "NULL",
				Tracks:    "enough",
			},
			Description:  "cooL?",
			SongDuration: timestamppb.Now()},
		musicpb.Music{SongName: "Idk",
			AuthorInfo:  &musicpb.Author{},
			Description: "someone", SongDuration: timestamppb.Now()},
	}

	stream, err := c.UpdateMusic(context.Background())
	if err != nil {
		log.Fatalf("Error while calling UpdateMusic RPC: %v\n", err)
	}

	for _, r := range req {
		stream.Send(&musicpb.UpdateMusicRequest{Req: &r})

	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server: %v\n", err)
	}
	fmt.Printf("%v\n", res)
}
