package main

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"

	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"

	"github.com/ozoncp/ocp-experience-api/config"
	"github.com/ozoncp/ocp-experience-api/internal/api"
	"github.com/ozoncp/ocp-experience-api/internal/db"
	"github.com/ozoncp/ocp-experience-api/internal/repo"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"

	sql "github.com/jmoiron/sqlx"
	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

const (
	apiKafkaTopic = "ocp_experience_events"
)

// creates kafka producer from config
func createKafkaProducer(config *config.Configuration) prod.Producer {
	brokers := config.KafkaEndpoint

	cfg := sarama.NewConfig()
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string { brokers }, cfg)

	if err != nil {
		log.Panic().Msgf("failed to connect to Kafka brokers: %v", err)
	}

	return prod.NewProducer(apiKafkaTopic, producer)
}

//
func initTracing(config *config.Configuration) {
	configuration := jaegerconfig.Configuration{
		ServiceName: "ocp-experience-api",
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LocalAgentHostPort: config.JaegerEndpoint,
			LogSpans:           true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := jaegermetrics.NullFactory

	tracer, _, err := configuration.NewTracer(
		jaegerconfig.Logger(jLogger),
		jaegerconfig.Metrics(jMetricsFactory),
	)

	if err != nil {
		log.Panic().Msgf("failed to initialize jaeger: %v", err)
	}

	opentracing.SetGlobalTracer(tracer)
}

func createExperienceApi(config *config.Configuration) *api.ExperienceAPI {
	database := db.Connect(config.ExperienceDNS)

	repo := repo.NewRepo(database)
	prom := metrics.NewMetricsReporter()
	producer := createKafkaProducer(config)
	tracer := opentracing.GlobalTracer()

	return api.NewRequestApi(repo, config.ExperienceBatchSize, prom, producer, tracer)
}

func run(database *sql.DB, config *config.Configuration) error {
	listen, err := net.Listen("tcp", ":" + strconv.FormatUint(config.GRPCPort, 10))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	desc.RegisterOcpExperienceApiServer(server, api.NewExperienceApi(repo.NewRepo(database)))

	serverErr := server.Serve(listen)

	if serverErr != nil {
		log.Fatalf("failed to serve: %v", serverErr)
	}

	return nil
}

func runJSON(config *config.Configuration) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpExperienceApiHandlerFromEndpoint(ctx, mux, config.GRPCServerEndpoint, opts)

	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(config.HTTPServerEndpoint, mux)

	if err != nil {
		panic(err)
	}
}

func main() {
	config := config.GetConfiguration("config.json")

	database := db.Connect(config.ExperienceDNS)
	defer database.Close()

	go runJSON(config)
	err := run(database, config)

	if err != nil {
		log.Fatal(err)
	}
}
