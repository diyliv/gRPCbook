package server

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	authpb "github.com/diyliv/grpc/5security/5auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *server) Add(ctx context.Context, r *authpb.AddReq) (*authpb.AddResp, error) {
	return &authpb.AddResp{
		Sum: int32(r.A + r.B),
	}, nil
}

type server struct{}

func NewServer() *server {
	return &server{}
}

var (
	crtFile = "../certs/leaf.pem"
	keyFile = "../certs/leaf.key"
)

func (s *server) Run() {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Error while parsing public/private key: %v\n", err)
	}
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(ensureValidBasicCreds),
	}
	serv := grpc.NewServer(opts...)
	authpb.RegisterSumServiceServer(serv, &server{})

	log.Printf("Starting gRPC server on port :50051")
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Error while listening: %v\n", err)
	}

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

func valid(authorization []string) bool {
	//функция для проверки являются ли данные, которые посылает юзер верными
	if len(authorization) < 1 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Basic ")
	fmt.Println(authorization)
	fmt.Println(token)
	return token == base64.StdEncoding.EncodeToString([]byte("admin:admin"))
}

func ensureValidBasicCreds(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Missing metadata")
	}

	if !valid(md["authorization"]) {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	return handler(ctx, req)
}
