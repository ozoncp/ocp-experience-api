package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/ozoncp/ocp-experience-api/internal/db"
	"github.com/ozoncp/ocp-experience-api/internal/repo"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-experience-api/internal/api"
	"google.golang.org/grpc"

	sql "github.com/jmoiron/sqlx"
	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

const (
	grpcPort           = ":82"
	grpcServerEndpoint = "localhost:82"
)

func mustGetEnvVar(name string) string {
	envVal := os.Getenv(name)

	if envVal == "" {
		panic(name + "is not set")
	}

	return envVal
}

func run(database *sql.DB) error {
	listen, err := net.Listen("tcp", grpcPort)

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

func runJSON() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpExperienceApiHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)

	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8082", mux)

	if err != nil {
		panic(err)
	}
}

func main() {
	dsn := mustGetEnvVar("OCP_EXPERIENCE_DSN")
	database := db.Connect(dsn)
	defer database.Close()

	go runJSON()
	err := run(database)

	if err != nil {
		log.Fatal(err)
	}
}
