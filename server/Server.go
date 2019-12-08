package server

import (
	context "context"
	"log"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	schmokinHTTP "github.com/reaandrew/schmokin/infrastructure/http"
	"github.com/reaandrew/schmokin/service"
	"github.com/reaandrew/schmokin/utils"
	grpc "google.golang.org/grpc"
)

var server *grpc.Server

type schmokinRemoteService struct {
}

func (s *schmokinRemoteService) Run(ctx context.Context, in *SchmokinRequest) (*SchmokinResponse, error) {
	schmokinService := service.NewSchmokinServiceBuilder().
		SetClient(schmokinHTTP.NewDefaultClient()).
		SetIterations(int(in.Iterations)).
		SetRandom(in.Random).
		SetTimer(utils.NewDefaultTimer()).
		SetWorkers(int(in.WorkerCount)).
		Build()

	result := schmokinService.Execute(in.Lines)

	response := &SchmokinResponse{
		Transactions:           int32(result.Transactions),
		Availability:           result.Availability,
		ElapsedTime:            int64(result.ElapsedTime),
		AverageResponseTime:    result.AverageResponseTime,
		ConcurrencyRate:        result.ConcurrencyRate,
		DataReceiveRate:        result.DataReceiveRate,
		DataSendRate:           result.DataSendRate,
		FailedTransactions:     result.FailedTransactions,
		LongestTransaction:     result.LongestTransaction,
		ShortestTransaction:    result.ShortestTransaction,
		SuccessfulTransactions: result.SuccessfulTransactions,
		TotalBytesReceived:     int32(result.TotalBytesReceived),
		TotalBytesSent:         int32(result.TotalBytesSent),
		TransactionRate:        result.TransactionRate,
	}
	return response, nil
}

func (s *schmokinRemoteService) Ping(ctx context.Context, in *empty.Empty) (*PingResponse, error) {
	return &PingResponse{
		Healthy: true,
	}, nil
}

func (s *schmokinRemoteService) Kill(ctx context.Context, in *empty.Empty) (*KillResponse, error) {
	server.Stop()
	return &KillResponse{
		Killed: true,
	}, nil
}

func StartServer(address string) {
	lis, err := net.Listen("tcp", address)
	log.Println("Server starting on " + address)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	server = grpc.NewServer()
	RegisterSchmokinServiceServer(server, &schmokinRemoteService{})

	if err := server.Serve(lis); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to start server!"))
	}
}
