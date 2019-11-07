package service_test

import (
	"context"
	"errors"

	. "github.com/sjeandeaux/ori/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HealthChecker", func() {
	Describe("Check", func() {
		Context("With a healthcheck without function", func() {
			It("should return UNKNOWN", func() {
				h := HealthChecker{}

				actual, err := h.Check(context.TODO(), nil)
				Ω(actual).Should(Equal(&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_UNKNOWN}))
				Ω(err).Should(BeNil())
			})
		})

		Context("With a healthcheck with function which returns nil as error", func() {
			It("should return SERVING", func() {
				h := HealthChecker{
					HealthCheck: func() error { return nil },
				}

				actual, err := h.Check(context.TODO(), nil)
				Ω(actual).Should(Equal(&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}))
				Ω(err).Should(BeNil())
			})
		})

		Context("With a healthcheck with function which returns not nil as error", func() {
			It("should return NOT_SERVING", func() {
				expectedErr := errors.New("booum")
				h := HealthChecker{
					HealthCheck: func() error { return expectedErr },
				}
				actual, err := h.Check(context.TODO(), nil)
				Ω(actual).Should(Equal(&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING}))
				Ω(err).Should(Equal(expectedErr))
			})
		})
	})

	Describe("Watch", func() {
		Context("With any healthcheck", func() {
			It("should return an error", func() {
				h := HealthChecker{
					HealthCheck: func() error { return nil },
				}
				err := h.Watch(nil, nil)
				Ω(err).Should(Equal(status.Errorf(codes.Unimplemented, "not implemented")))
			})
		})
	})
})
