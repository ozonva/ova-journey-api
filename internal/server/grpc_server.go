package server

import (
	"github.com/jmoiron/sqlx"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ozonva/ova-journey-api/internal/api"
	"github.com/ozonva/ova-journey-api/internal/config"
	"github.com/ozonva/ova-journey-api/internal/repo"
	desc "github.com/ozonva/ova-journey-api/pkg/ova-journey-api"
)

// GrpcServer - represents simple gRPC server wrapper
type GrpcServer struct {
	configuration *config.EndpointConfiguration
	db            *sqlx.DB
	server        *grpc.Server
	errChan       chan<- error
	chunkSize     int
}

// NewGrpcServer - creates new GrpcServer with configuration endpoint
//
// and output channel to signalize about critical errors
func NewGrpcServer(configuration *config.EndpointConfiguration, db *sqlx.DB, chunkSize int, errChan chan<- error) *GrpcServer {
	return &GrpcServer{
		configuration: configuration,
		db:            db,
		errChan:       errChan,
		chunkSize:     chunkSize,
	}
}

// Start - start GrpcServer
func (s *GrpcServer) Start() {
	endpointAddress := s.configuration.GetEndpointAddress()
	listen, err := net.Listen("tcp4", endpointAddress)
	if err != nil {
		log.Err(err).Msg("GRPC server: failed to listen")
		s.errChan <- err
	}

	repository := repo.NewRepo(s.db)

	s.server = grpc.NewServer()
	desc.RegisterJourneyApiV1Server(s.server, api.NewJourneyAPI(repository, s.chunkSize))

	go func() {
		log.Debug().Msg("GRPC server: starting")
		if err := s.server.Serve(listen); err != nil {
			log.Err(err).Msg("GRPC server: failed to serve")
			s.errChan <- err
		}
	}()
}

// Stop - graceful stop GrpcServer
func (s *GrpcServer) Stop() {
	s.server.GracefulStop()
}
