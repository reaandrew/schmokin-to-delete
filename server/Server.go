package server

import (
	context "context"
	"log"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	surgeHTTP "github.com/reaandrew/surge/infrastructure/http"
	"github.com/reaandrew/surge/service"
	"github.com/reaandrew/surge/utils"
	grpc "google.golang.org/grpc"
)

var server *grpc.Server

type surgeRemoteService struct {
}

func (s *surgeRemoteService) Run(ctx context.Context, in *SurgeRequest) (*SurgeResponse, error) {
	surgeService := service.NewSurgeServiceBuilder().
		SetClient(surgeHTTP.NewDefaultClient()).
		SetIterations(int(in.Iterations)).
		SetRandom(in.Random).
		SetTimer(utils.NewDefaultTimer()).
		SetWorkers(int(in.WorkerCount)).
		Build()

	result := surgeService.Execute(in.Lines)

	response := &SurgeResponse{
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

func (s *surgeRemoteService) Ping(ctx context.Context, in *empty.Empty) (*PingResponse, error) {
	return &PingResponse{
		Healthy: true,
	}, nil
}

func (s *surgeRemoteService) Kill(ctx context.Context, in *empty.Empty) (*KillResponse, error) {
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
	RegisterSurgeServiceServer(server, &surgeRemoteService{})

	if err := server.Serve(lis); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to start server!"))
	}
}
