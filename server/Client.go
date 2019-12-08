package server

import (
	"log"

	grpc "google.golang.org/grpc"
)

func CreateClient(endpoint string) SchmokinServiceClient {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := NewSchmokinServiceClient(conn)
	return c
}
