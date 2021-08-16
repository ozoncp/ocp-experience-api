package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-experience-api/internal/api"
	"google.golang.org/grpc"

	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

const (
	grpcPort           = ":82"
	grpcServerEndpoint = "localhost:82"
)

func run() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	desc.RegisterOcpExperienceApiServer(server, api.NewExperienceApi())

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
	go runJSON()
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}
