package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "github.com/shtsukada/k8s-observability-app/gen/proto"
	"google.golang.org/grpc"
)

var (
	grpcRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "grpc_server_requests_total",
			Help: "Total number of gRPC requests received",
		},
	)
)

type server struct {
	pb.UnimplementedEchoServiceServer
}

func (s *server) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	grpcRequests.Inc()
	log.Printf("Received gRPC request: %s", req.Message)
	return &pb.EchoResponse{Message: "Echo: " + req.Message}, nil
}

func main() {
	prometheus.MustRegister(grpcRequests)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterEchoServiceServer(s, &server{})
		log.Println("gRPC server listening on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics endpoint at :8080/metrics")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	log.Println("All services running")
	select {}
}
