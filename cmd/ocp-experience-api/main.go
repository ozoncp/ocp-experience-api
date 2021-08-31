package main

import (
	"context"
	"encoding/json"
	"go.uber.org/atomic"

	"github.com/ozoncp/ocp-experience-api/internal/metrics"
	"github.com/ozoncp/ocp-experience-api/internal/producer"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"net"
	"net/http"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"

	"github.com/uber/jaeger-client-go"

	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	jaegerMetrics "github.com/uber/jaeger-lib/metrics"

	"github.com/ozoncp/ocp-experience-api/config"
	"github.com/ozoncp/ocp-experience-api/internal/api"
	"github.com/ozoncp/ocp-experience-api/internal/db"
	"github.com/ozoncp/ocp-experience-api/internal/repo"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"

	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

var isServiceReady atomic.Bool

const (
	apiKafkaTopic = "ocp_experience_events"
)

// creates kafka producer from config
func createKafkaProducer(config *config.Configuration) producer.Producer {
	brokers := config.KafkaEndpoint

	cfg := sarama.NewConfig()
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true

	prod, err := sarama.NewSyncProducer([]string { brokers }, cfg)

	if err != nil {
		log.Panic().Msgf("failed to connect to Kafka brokers: %v", err)
	}

	return producer.NewProducer(apiKafkaTopic, prod)
}

// inits opentracing
func initTracing(config *config.Configuration) {
	configuration := jaegerConfig.Configuration{
		ServiceName: "ocp-experience-api",
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort: config.JaegerEndpoint,
			LogSpans:           true,
		},
	}

	jLogger := jaegerLog.StdLogger
	jMetricsFactory := jaegerMetrics.NullFactory

	tracer, _, err := configuration.NewTracer(
		jaegerConfig.Logger(jLogger),
		jaegerConfig.Metrics(jMetricsFactory),
	)

	if err != nil {
		log.Panic().Msgf("failed to initialize jaeger: %v", err)
	}

	opentracing.SetGlobalTracer(tracer)
}

// builds experience API service
func createExperienceApi(config *config.Configuration) *api.ExperienceAPI {
	database := db.Connect(config.ExperienceDNS)

	repo := repo.NewRepo(database)
	prom := metrics.NewReporter()
	producer := createKafkaProducer(config)
	tracer := opentracing.GlobalTracer()

	return api.NewExperienceApi(repo, config.ExperienceBatchSize, prom, producer, tracer)
}

func run(config *config.Configuration) error {
	listen, err := net.Listen("tcp", ":" + strconv.FormatUint(config.GRPCPort, 10))

	if err != nil {
		log.Panic().Msgf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	experienceApi := createExperienceApi(config)

	desc.RegisterOcpExperienceApiServer(server, experienceApi)

	isServiceReady.Store(true)
	serverErr := server.Serve(listen)

	if serverErr != nil {
		log.Fatal().Msgf("failed to serve: %v", serverErr)
		return serverErr
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

// runs metrics
func runUtilityServer() {
	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/live", http.HandlerFunc(liveHandler))
	http.Handle("/health", http.HandlerFunc(healthHandler))
	http.Handle("/ready", readyHandler(&isServiceReady))

	if err := http.ListenAndServe(":9100", nil); err != nil {
		log.Panic().Msgf("http endpoint failed: %v", err)
	}
}

func main() {
	config := config.GetConfiguration("config.json")

	initTracing(config)

	go runJSON(config)
	go runUtilityServer()

	err := run(config)

	if err != nil {
		log.Fatal().Msgf("Error %e", err)
	}
}

// liveHandler is a simple HTTP handler function which writes a response
func liveHandler(w http.ResponseWriter, _ *http.Request) {
	body, err := json.Marshal("This service is alive")

	if err != nil {
		log.Error().Err(err).Msgf("Could not encode info data: %v", err)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)

	if err != nil {
		log.Error().Err(err).Msgf("Could not write body: %v", err)
		return
	}
}

// healthHandler is a Liveness probe
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readyHandler is a Readiness probe
func readyHandler(isReady *atomic.Bool) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load() {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
