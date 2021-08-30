package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ozonva/ova-journey-api/internal/config"
	desc "github.com/ozonva/ova-journey-api/pkg/ova-journey-api"
)

// GatewayServer - represents simple http server wrapper for gateway to gRPC service
type GatewayServer struct {
	gatewayConfiguration      *config.EndpointConfiguration
	grpcEndpointConfiguration *config.EndpointConfiguration
	httpServer                *http.Server
	errChan                   chan<- error
	wg                        *sync.WaitGroup
}

// NewGatewayServer - creates new GatewayServer with configuration parameters
//
// and output channel to signalize about critical errors
func NewGatewayServer(gatewayConfiguration, grpcEndpointConfiguration *config.EndpointConfiguration, errChan chan<- error) *GatewayServer {
	return &GatewayServer{
		gatewayConfiguration:      gatewayConfiguration,
		grpcEndpointConfiguration: grpcEndpointConfiguration,
		errChan:                   errChan,
	}
}

// Start - start GatewayServer
func (s *GatewayServer) Start() {
	gatewayAddress := s.gatewayConfiguration.GetEndpointAddress()
	grpcAddress := s.grpcEndpointConfiguration.GetEndpointAddress()

	mux := http.NewServeMux()
	gatewayMux := runtime.NewServeMux()
	mux.Handle("/", gatewayMux)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))

	s.httpServer = &(http.Server{Addr: gatewayAddress, Handler: mux})

	s.wg = &sync.WaitGroup{}

	go func() {
		opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := desc.RegisterJourneyApiV1HandlerFromEndpoint(ctx, gatewayMux, grpcAddress, opts)
		if err != nil {
			log.Err(err).Msg("Gateway server: failed to register handler to GRPC")
			s.errChan <- err
		}

		log.Debug().Msg("Gateway server: starting")
		s.wg.Add(1)
		err = s.httpServer.ListenAndServe()

		if err == http.ErrServerClosed {
			s.wg.Done()
			return
		}

		if err != nil {
			log.Err(err).Msg("Gateway server: failed to serve")
			s.errChan <- err
		}
	}()
}

// Stop - graceful stop GatewayServer with waiting of http server is stopped
func (s *GatewayServer) Stop() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Err(err).Msg("Gateway server: shutdown failed")
		s.errChan <- err
	}
	s.wg.Wait()
}
