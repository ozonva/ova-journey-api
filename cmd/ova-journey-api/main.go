package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	db, grpc, gateway := startApp(configuration, errChan)

	cu.WatchConfigurationFile(func(configuration config.Configuration) {
		configChan <- configuration
	})

	for {
		select {
		case c := <-configChan:
			log.Info().Msg("Restart service after changing configuration")
			stopApp(db, grpc, gateway)
			db, grpc, gateway = startApp(&c, errChan)
			log.Debug().Msg("Restart service success")
		case err := <-errChan:
			log.Err(err).Msg("Internal server error")
			stopApp(db, grpc, gateway)
			return
		case <-quit:
			log.Info().Msg("Shutdown service")
			stopApp(db, grpc, gateway)
			return
		}
	}
}

func startApp(c *config.Configuration, errChan chan<- error) (*sqlx.DB, *server.GrpcServer, *server.GatewayServer) {
	db, err := createDb(c.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot establish connection to database")
		return nil, nil, nil
	}
	grpcServer := server.NewGrpcServer(c.GRPC, db, errChan)
	gatewayServer := server.NewGatewayServer(c.Gateway, c.GRPC, errChan)
	grpcServer.Start()
	gatewayServer.Start()
	return db, grpcServer, gatewayServer
}

func stopApp(database *sqlx.DB, grpcServer *server.GrpcServer, gatewayServer *server.GatewayServer) {
	gatewayServer.Stop()
	grpcServer.Stop()
	if err := database.Close(); err != nil {
		log.Fatal().Err(err).Msg("Database close error")
	}
}

func createDb(configuration *config.DatabaseConfiguration) (*sqlx.DB, error) {
	db, err := sqlx.Open(configuration.Driver, configuration.GetDataSourceName())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
