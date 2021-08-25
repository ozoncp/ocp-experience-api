package config

import "github.com/tkanos/gonfig"

// config default values
const (
	experienceDNS = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	grpcPort = 82
	grpcServerEndpoint = "localhost:82"
	httpServerEndpoint = "localhost:8082"
)

// Configuration describes app config
type Configuration struct {
	ExperienceDNS string	// PostgresSQL DNS endpoint
	GRPCPort uint64
	GRPCServerEndpoint string
	HTTPServerEndpoint string
}

// GetConfiguration reads config file and returns config as struct
// If fileName does not exist, then returns default config
func GetConfiguration(fileName string)  *Configuration {
	var config Configuration
	err := gonfig.GetConf(fileName, &config)

	if err != nil {
		fillDefaultConfig(&config)
	}

	return &config
}

// fills Configuration with default values
func fillDefaultConfig(config *Configuration) {
	config.ExperienceDNS = experienceDNS
	config.GRPCPort = grpcPort
	config.GRPCServerEndpoint = grpcServerEndpoint
	config.HTTPServerEndpoint = httpServerEndpoint
}
