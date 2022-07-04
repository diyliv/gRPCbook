package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	musicpb "github.com/diyliv/grpc/startWorkWithgRPC/ex/proto/music"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type controller struct { // (0)
	musicInfo map[string]musicpb.Music // (1)
	mu        sync.Mutex
}

func NewController(musicInfo map[string]musicpb.Music) *controller {
	return &controller{
		musicInfo: musicInfo,
	}
}

func (c *controller) AddMusic(ctx context.Context, r *musicpb.MusicRequest) (*musicpb.MusicID, error) { // (2)
	id, err := uuid.NewV4()
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while generating unique ID")
	}

	c.mu.Lock()
	c.musicInfo[id.String()] = musicpb.Music{
		SongName: r.Add.GetSongName(),
		AuthorInfo: &musicpb.Author{
			Name:      r.Add.GetAuthorInfo().Name,
			Age:       r.Add.GetAuthorInfo().Age,
			Followers: r.Add.GetAuthorInfo().Followers,
			Tracks:    r.Add.GetAuthorInfo().Tracks,
		},
		Description:  r.Add.GetDescription(),
		SongDuration: r.Add.GetSongDuration(),
	}
	c.mu.Unlock()

	return &musicpb.MusicID{Id: id.String()}, nil
}

func (c *controller) GetMusic(ctx context.Context, r *musicpb.MusicID) (*musicpb.MusicResponse, error) { // (3)
	id := r.GetId()

	c.mu.Lock()
	val, ok := c.musicInfo[id]
	if !ok {
		return nil, status.Error(codes.NotFound, "No such track")
	}
	c.mu.Unlock()
	return &musicpb.MusicResponse{Resp: &val}, nil
}

func (c *controller) SearchMusic(req *wrapperspb.StringValue, stream musicpb.MusicService_SearchMusicServer) error {
	name := req.GetValue()
	fmt.Println(name)

	peformers := make([]musicpb.Music, 0)
	peformers = append(peformers, musicpb.Music{
		SongName: "Lofi for studying 10hours",
		AuthorInfo: &musicpb.Author{
			Name:      "YT channel",
			Age:       "NULL",
			Followers: "1k",
			Tracks:    "NULL",
		},
		Description:  "chill lofi songs for studying/sleeping/relaxing",
		SongDuration: timestamppb.Now(),
	},
		musicpb.Music{
			SongName: "Rain 10 hours",
			AuthorInfo: &musicpb.Author{
				Name:      "YT channel",
				Age:       "NULL",
				Followers: "2k",
				Tracks:    "NULL",
			},
			Description:  "rain 10 hours for sleeping",
			SongDuration: timestamppb.Now()},
		musicpb.Music{SongName: "Forest sounds",
			AuthorInfo: &musicpb.Author{
				Name:      "YT channel",
				Age:       "NULL",
				Followers: "3k",
				Tracks:    "NULL",
			},
			Description:  "forest sounds",
			SongDuration: timestamppb.Now()})

	for i := 0; i < len(peformers); i++ {
		res := musicpb.SearchMusicResponse{
			Search: &peformers[i],
		}
		if err := stream.Send(&res); err != nil {
			log.Fatalf("Error while sending stream response: %v\n", err)
		}
	}
	return nil
}

func (c *controller) UpdateMusic(stream musicpb.MusicService_UpdateMusicServer) error {
	for {
		req, err := stream.Recv() // (1)
		if err == io.EOF {        // (2)
			return stream.SendAndClose(&wrapperspb.StringValue{
				Value: "Tracks were updated",
			})
		}
		fmt.Println(req)
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}
	}
}

func (c *controller) UploadMusic(stream musicpb.MusicService_UploadMusicServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while receiving data: %v\n", err)
		}

		if err := stream.Send(&musicpb.UploadMusicResponse{
			Resp: "Uploaded " + req.Req.TrackName,
		}); err != nil {
			log.Fatalf("Error while sending stream response: %v\n", err)
		}
	}
}
