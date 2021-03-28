package main

import (
	"context"
	pb "fibServer/proto"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var port string

type WorkerServer struct {
	pb.UnimplementedFibWorkerServer
}

// given a index number returns its fib number
func (ws *WorkerServer) GetFibNumber(ctx context.Context, req *pb.FibRequest) (*pb.FibReply, error) {
	if req.Num < 0 {
		return nil, fmt.Errorf("can't calculate negative index of fibonachi numbers")
	}
	return &pb.FibReply{Num: calculateFib(req.Num)}, nil
}

func calculateFib(index int32) int32 {
	if index <= 2 {
		return 1
	}
	return calculateFib(index-1) + calculateFib(index-2)
}

func main() {
	flag.StringVar(&port, "port", "localhost:5544", "provide IP addr and port for running the server")
	flag.Parse()
	fmt.Printf("Running the server with: %s\n", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption // options for the server

	grpcServer := grpc.NewServer(opts...) // create the grpc server instance
	fibWorkServer := WorkerServer{}       // create our server instance

	pb.RegisterFibWorkerServer(grpcServer, &fibWorkServer)
	grpcServer.Serve(lis)
}
