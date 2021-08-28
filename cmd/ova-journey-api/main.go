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

	errGrpcChan := make(chan error, 1)
	errGatewayChan := make(chan error, 1)
	configChan := make(chan config.Configuration)

	cu := config.NewConfigurationUpdater(ConfigUpdatePeriod, ConfigFile)
	configuration := cu.GetConfiguration()

	grpcServer := server.NewGrpcServer(configuration.GRPC, errGrpcChan)
	gatewayServer := server.NewGatewayServer(configuration.Gateway, configuration.GRPC, errGatewayChan)
	grpcServer.Start()
	gatewayServer.Start()

	cu.WatchConfigurationFile(func(configuration config.Configuration) {
		configChan <- configuration
	})

	for {
		select {
		case configuration := <-configChan:
			log.Info().Msg("Restart servers after changing configuration")
			gatewayServer.Stop()
			grpcServer.Stop()
			gatewayServer.UpdateConfiguration(configuration.Gateway, configuration.GRPC)
			grpcServer.UpdateConfiguration(configuration.GRPC)
			grpcServer.Start()
			gatewayServer.Start()
			log.Debug().Msg("Restart servers success")
		case err := <-errGrpcChan:
			log.Err(err).Msg("GRPC server error")
			gatewayServer.Stop()
			grpcServer.Stop()
			return
		case err := <-errGatewayChan:
			log.Err(err).Msg("Gateway server error")
			gatewayServer.Stop()
			grpcServer.Stop()
			return
		case <-quit:
			log.Info().Msg("Shutdown service")
			gatewayServer.Stop()
			grpcServer.Stop()
			return
		}
	}
}
