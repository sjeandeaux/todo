package http

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RunServer runs the http server with
func RunServer(ctx context.Context, host, httpPort string) (int, error) {

	lis, err := net.Listen("tcp", net.JoinHostPort(host, httpPort))
	if err != nil {
		return -1, err
	}

	ctx, cancel := context.WithCancel(ctx)
	go func(list net.Listener, ctx context.Context, cancel context.CancelFunc) {

		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		s := &http.Server{
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      mux,
		}

		if err := s.Serve(lis); err != nil {
			log.Println(err) //TODO manage the error
			cancel()
		}

	}(lis, ctx, cancel)

	return lis.Addr().(*net.TCPAddr).Port, nil
}
