package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

//HealthChecker the health checker
type HealthChecker struct {
	HealthCheck func() error
}

// Check implements the health check interface, which directly returns to health status. There are also more complex health check strategies, such as returning based on server load.
func (h *HealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.Info("received health check")
	if h.HealthCheck == nil {
		log.Error("HealthCheck function is nil")
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_UNKNOWN}, nil
	}

	if err := h.HealthCheck(); err != nil {
		log.Error(err)
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING}, err
	}

	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

// Watch not yet implemented
func (h *HealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "not implemented")
}
