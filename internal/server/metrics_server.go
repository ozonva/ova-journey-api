package server

import (
	"context"
	"github.com/ozonva/ova-journey-api/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
)

// MetricsServer - represents simple http server wrapper for prometheus metrics
type MetricsServer struct {
	prometheusConfiguration *config.PrometheusConfiguration
	httpServer              *http.Server
	wg                      *sync.WaitGroup
}

// NewMetricsServer - creates new MetricsServer with configuration parameters
func NewMetricsServer(prometheusConfiguration *config.PrometheusConfiguration) *MetricsServer {
	return &MetricsServer{
		prometheusConfiguration: prometheusConfiguration,
	}
}

// Start - start MetricsServer
func (s *MetricsServer) Start() {
	mux := http.NewServeMux()
	mux.Handle(s.prometheusConfiguration.Path, promhttp.Handler())

	s.httpServer = &http.Server{
		Addr:    s.prometheusConfiguration.GetEndpointAddress(),
		Handler: mux,
	}

	s.wg = &sync.WaitGroup{}

	go func() {
		log.Debug().Msg("Metrics server: starting")
		s.wg.Add(1)
		err := s.httpServer.ListenAndServe()

		if err == http.ErrServerClosed {
			s.wg.Done()
			return
		}

		if err != nil {
			log.Err(err).Msg("Metrics server: failed to serve")
		}
	}()
}

// Stop - graceful stop MetricsServer with waiting of http server is stopped
func (s *MetricsServer) Stop() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Err(err).Msg("Metrics server: shutdown failed")
	}
	s.wg.Wait()
}
