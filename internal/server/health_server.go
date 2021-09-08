package server

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/ozonva/ova-journey-api/internal/config"
	"github.com/ozonva/ova-journey-api/internal/kafka"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
)

// HealthServer - represents simple http server for checking health
type HealthServer struct {
	healthcheckConfiguration *config.HealthCheckConfiguration
	httpServer               *http.Server
	wg                       *sync.WaitGroup
	producer                 kafka.Producer
	db                       *sqlx.DB
}

// NewHealthServer - creates new HealthServer with configuration parameters
func NewHealthServer(healthcheckConfiguration *config.HealthCheckConfiguration, producer kafka.Producer, db *sqlx.DB) *HealthServer {
	return &HealthServer{
		healthcheckConfiguration: healthcheckConfiguration,
		producer:                 producer,
		db:                       db,
	}
}

// Start - start HealthServer
func (s *HealthServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc(s.healthcheckConfiguration.Path, s.healthHandler())

	s.httpServer = &http.Server{
		Addr:    s.healthcheckConfiguration.GetEndpointAddress(),
		Handler: mux,
	}

	s.wg = &sync.WaitGroup{}

	go func() {
		log.Debug().Msg("HealthServer: starting")
		s.wg.Add(1)
		err := s.httpServer.ListenAndServe()

		if err == http.ErrServerClosed {
			s.wg.Done()
			return
		}

		if err != nil {
			log.Err(err).Msg("HealthServer: failed to serve")
		}
	}()
}

// Stop - graceful stop HealthServer with waiting of http server is stopped
func (s *HealthServer) Stop() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Err(err).Msg("HealthServer: shutdown failed")
	}
	s.wg.Wait()
}

func (s *HealthServer) healthHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		servicesState := make([]string, 0, 2)
		servicesState = append(servicesState, s.checkDbHealth())
		servicesState = append(servicesState, s.checkKafkaHealth())

		err := json.NewEncoder(w).Encode(servicesState)
		if err != nil {
			log.Error().Err(err).Msg("Error in encoding internal services state")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s *HealthServer) checkDbHealth() string {
	_, err := s.db.Query("SELECT 1")
	if err != nil {
		return "DB: Failed"
	}
	return "DB: OK"
}

func (s *HealthServer) checkKafkaHealth() string {
	err := s.producer.Send(kafka.Message{MessageType: kafka.Ping, Value: "1"})
	if err != nil {
		return "Kafka: Failed"
	}
	return "Kafka: OK"
}
