// Journeys API (gRPC + JSON gateway) for Ozon Voice Assistant project
//
// Ozon Voice Assistant Project Documentation is here https://github.com/ozonva/docs
package main

import (
	"io"
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
	"github.com/ozonva/ova-journey-api/internal/kafka"
	"github.com/ozonva/ova-journey-api/internal/metrics"
	"github.com/ozonva/ova-journey-api/internal/server"
	"github.com/ozonva/ova-journey-api/internal/tracer"
)

//ConfigFile - application configuration file path
const ConfigFile = "config/config.yaml"

// ConfigUpdatePeriod - time duration between checking updates in configuration file
const ConfigUpdatePeriod = 5 * time.Second

var (
	db            *sqlx.DB
	grpc          *server.GrpcServer
	gateway       *server.GatewayServer
	healthChecker *server.HealthServer
	tracerCloser  io.Closer
	producer      kafka.Producer
	metricServer  *server.MetricsServer
	metric        metrics.Metrics
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	errChan := make(chan error, 1)
	configChan := make(chan config.Configuration)

	cu := config.NewConfigurationUpdater(ConfigUpdatePeriod, ConfigFile)
	configuration := cu.GetConfiguration()
	log.Info().Str("version", configuration.Project.Version).Msg("Starting ova-journey-api")

	metric = metrics.NewMetrics("ova_journey_api", "gRPC_server")

	startApp(configuration, errChan)

	cu.WatchConfigurationFile(func(configuration config.Configuration) {
		configChan <- configuration
	})

	for {
		select {
		case c := <-configChan:
			log.Info().Msg("Restart service after changing configuration")
			stopApp()
			startApp(&c, errChan)
			log.Debug().Msg("Restart service success")
		case err := <-errChan:
			log.Err(err).Msg("Internal server error")
			stopApp()
			return
		case <-quit:
			log.Info().Msg("Shutdown service")
			stopApp()
			return
		}
	}
}

func startApp(c *config.Configuration, errChan chan<- error) {
	var err error

	tracerCloser = tracer.InitTracer(c.Project.Name, c.Jaeger)

	producer, err = kafka.NewProducer(c.Kafka)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Kafka producer")
	}

	db, err = createDb(c.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot establish connection to database")
	}

	healthChecker = server.NewHealthServer(c.HealthCheck, producer, db)
	metricServer = server.NewMetricsServer(c.Prometheus)
	grpc = server.NewGrpcServer(c.GRPC, producer, db, metric, c.ChunkSize, errChan)
	gateway = server.NewGatewayServer(c.Gateway, c.GRPC, errChan)

	healthChecker.Start()
	metricServer.Start()
	grpc.Start()
	gateway.Start()
}

func stopApp() {
	gateway.Stop()
	grpc.Stop()
	metricServer.Stop()
	healthChecker.Stop()
	if err := db.Close(); err != nil {
		log.Fatal().Err(err).Msg("Database close error")
	}

	if err := producer.Close(); err != nil {
		log.Fatal().Err(err).Msg("Kafka producer close error")
	}

	if err := tracerCloser.Close(); err != nil {
		log.Fatal().Err(err).Msg("Tracer close error")
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
