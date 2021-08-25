package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/ozoncp/ocp-experience-api/config"
	"github.com/ozoncp/ocp-experience-api/internal/api"
	"github.com/ozoncp/ocp-experience-api/internal/db"
	"github.com/ozoncp/ocp-experience-api/internal/repo"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"

	sql "github.com/jmoiron/sqlx"
	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

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
