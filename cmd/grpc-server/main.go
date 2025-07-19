package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	pb "github.com/shtsukada/k8s-observability-app/gen/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	logger       *zap.Logger
	tracer       trace.Tracer
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

	ctx, span := tracer.Start(ctx, "Echo")
	defer span.End()

	logger.Info("gRPC request received",
		zap.String("message", req.Message),
		zap.String("traceID", span.SpanContext().TraceID().String()),
		zap.String("spanID", span.SpanContext().SpanID().String()),
	)
	return &pb.EchoResponse{Message: "Echo: " + req.Message}, nil
}

func initLogger() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
}

func initTracer(ctx context.Context) (func(), error) {
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	tracer = otel.Tracer("grpc-app")
	return func() { _ = tp.Shutdown(ctx) }, nil
}

func main() {
	ctx := context.Background()

	initLogger()
	defer logger.Sync()

	shutdown, err := initTracer(ctx)
	if err != nil {
		logger.Fatal("failed to initialize tracer", zap.Error(err))
	}
	defer shutdown()

	prometheus.MustRegister(grpcRequests)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			logger.Fatal("failed to listen", zap.Error(err))
		}
		s := grpc.NewServer()
		pb.RegisterEchoServiceServer(s, &server{})
		logger.Info("gRPC server listening on :50051")
		if err := s.Serve(lis); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Info("Prometheus metrics endpoint at :8080/metrics")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger.Fatal("metrics server error", zap.Error(err))
		}
	}()

	logger.Info("gRPC + Metrics + Traces + Logs running")
	select {}
}
