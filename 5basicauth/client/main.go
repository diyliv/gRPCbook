package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	authpb "github.com/diyliv/grpc/5security/5auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	addr     = "localhost:50051"
	hostname = "localhost"
	crtFile  = "../certs/leaf.pem"
)

type basicAuth struct {
	username string
	password string
}

func (b basicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	auth := b.username + ":" + b.password

	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (b basicAuth) RequireTransportSecurity() bool {
	return true
}

func main() {
	creds, err := credentials.NewClientTLSFromFile(crtFile, hostname)
	if err != nil {
		log.Fatalf("Error while constrcuts TLS creds: %v\n", err)
	}

	auth := basicAuth{
		username: "admin",
		password: "admin",
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(auth), // interface
		grpc.WithTransportCredentials(creds),
	}

	cc, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("Error while dialing with gRPC server: %v\n", err)
	}

	defer cc.Close()

	c := authpb.NewSumServiceClient(cc)
	req, err := c.Add(context.Background(), &authpb.AddReq{
		A: int32(10),
		B: int32(5),
	})
	if err != nil {
		log.Fatalf("error while calling Add RPC: %v\n", err)
	}

	fmt.Println(req.Sum)
}
