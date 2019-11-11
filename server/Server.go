package server

import (
	context "context"
	"log"
	"net"

	"github.com/pkg/errors"
	grpc "google.golang.org/grpc"
)

const port = ":50051"

type surgeRemoteService struct {
}

func (s *surgeRemoteService) Run(ctx context.Context, in *SurgeRequest) (*SurgeResponse, error) {
	return nil, nil
}

func StartServer() *grpc.Server {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	server := grpc.NewServer()
	RegisterSurgeServiceServer(server, &surgeRemoteService{})

	if err := server.Serve(lis); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to start server!"))
	}

	return server
}
