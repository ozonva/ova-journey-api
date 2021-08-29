package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-journey-api/internal/config"
	"github.com/ozonva/ova-journey-api/internal/server"
)

//ConfigFile - application configuration file path
const ConfigFile = "config/config.yaml"

// ConfigUpdatePeriod - time duration between checking updates in configuration file
const ConfigUpdatePeriod = 5 * time.Second

func main() {
	log.Info().Msg("Hello, I'm ova-journey-api")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error, 1)
	configChan := make(chan config.Configuration)

	cu := config.NewConfigurationUpdater(ConfigUpdatePeriod, ConfigFile)
	configuration := cu.GetConfiguration()

	grpc, gateway := startApp(configuration, errChan)

	cu.WatchConfigurationFile(func(configuration config.Configuration) {
		configChan <- configuration
	})

	for {
		select {
		case c := <-configChan:
			log.Info().Msg("Restart service after changing configuration")
			stopApp(grpc, gateway)
			grpc, gateway = startApp(&c, errChan)
			log.Debug().Msg("Restart service success")
		case err := <-errChan:
			log.Err(err).Msg("Internal server error")
			stopApp(grpc, gateway)
			return
		case <-quit:
			log.Info().Msg("Shutdown service")
			stopApp(grpc, gateway)
			return
		}
	}
}

func startApp(c *config.Configuration, errChan chan<- error) (*server.GrpcServer, *server.GatewayServer) {
	grpcServer := server.NewGrpcServer(c.GRPC, errChan)
	gatewayServer := server.NewGatewayServer(c.Gateway, c.GRPC, errChan)
	grpcServer.Start()
	gatewayServer.Start()
	return grpcServer, gatewayServer
}

func stopApp(grpcServer *server.GrpcServer, gatewayServer *server.GatewayServer) {
	gatewayServer.Stop()
	grpcServer.Stop()
}
