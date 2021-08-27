# Ozon Code Platform Experience API

Experience API. The service accepts gRPC connections at port 82 and HTTP at 8082.

Supports:

- Create new experience
- MultiCreate new experiences
- Return experience information
- Remove experience
- Get experience list
- Update experience

### To build locally

- Install `protoc`. See instruction [here](https://grpc.io/docs/protoc-installation/)
- Build:

```shell
git clone https://github.com/ozoncp/ocp-experience-api.git
cd ocp-experience-api
make build
```
Built binary exists at `bin/ocp-experience-api` <br />
Create all tables - `make migrate` <br />
Run service - `docker compose up` <br />
Run tests - `make test` <br />

### To build and run with Docker

- Build docker image `docker build . -t ocp-experience-api`
- Run `docker run -p 82:82 ocp-experience-api`

### Configuration
Use `config.json` to configure instance, params are follows

- `ExperienceDNS`, by default is "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" - defines connection to PostgresSQL
- `GRPCPort`, by default is 82
- `GRPCServerEndpoint`, by default is "localhost:82"
- `HTTPServerEndpoint`, by default is "localhost:8082"
- `ExperienceBatchSize`, by default is 1000
- `KafkaEndpoint`, by default is "kafka:9094"
- `JaegerEndpoint`, by default is "jaeger:6831"
