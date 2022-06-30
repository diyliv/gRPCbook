package controllers

import (
	"context"
	"sync"

	musicpb "github.com/diyliv/grpc/startWorkWithgRPC/ex/proto/music"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
