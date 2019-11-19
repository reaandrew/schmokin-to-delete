package server

import (
	context "context"
	"log"
	"time"

	grpc "google.golang.org/grpc"
)

const address = "localhost:50051"

func CreateClient(endpoint string) SurgeServiceClient {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := NewSurgeServiceClient(conn)
	return c
}

func RunClient() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)

	}
	defer conn.Close()

	c := NewSurgeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.Run(ctx, &SurgeRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting")
}
