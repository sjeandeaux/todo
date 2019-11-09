package grpc

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"

	"github.com/sjeandeaux/todo/pkg/service"
	pb "github.com/sjeandeaux/todo/todo-grpc/v1"

	"google.golang.org/grpc/health/grpc_health_v1"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc/reflection"
)

// RunServer runs the grpc server on port port
func RunServer(ctx context.Context, host string, port string, server *service.ToDoServiceServer) (int, error) {
	const ProtoTCP = "tcp"
	lis, err := net.Listen(ProtoTCP, net.JoinHostPort(host, port))
	if err != nil {
		return -1, err
	}

	logrusEntry := log.NewEntry(log.StandardLogger())
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	ctx, cancel := context.WithCancel(ctx)
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_prometheus.StreamServerInterceptor,
			grpc_logrus.StreamServerInterceptor(logrusEntry),
		),
	)

	grpc_prometheus.Register(grpcServer)

	go func(ctx context.Context, cancel context.CancelFunc, grpcServer *grpc.Server, lis net.Listener) {
		pb.RegisterToDoServiceServer(grpcServer, server)
		healthCheck := &service.HealthChecker{
			HealthCheck: server.HealthChecher(),
		}
		grpc_health_v1.RegisterHealthServer(grpcServer, healthCheck)
		reflection.Register(grpcServer)
		if err := grpcServer.Serve(lis); err != nil {
			log.Println(err) //TODO manage the error
			cancel()
		}
	}(ctx, cancel, grpcServer, lis)

	go func(ctx context.Context, grpcServer *grpc.Server) {
		for {
			select {
			case <-ctx.Done():
				grpcServer.GracefulStop()
				return
			}
		}
	}(ctx, grpcServer)

	return lis.Addr().(*net.TCPAddr).Port, nil
}
