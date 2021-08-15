# Ozon Code Platform Experience API

Experience API. The service accepts gRPC connections at port 82 and HTTP at 8082.

Supports:

- Create new experience
- Return experience information
- Remove experience
- Get experience list

### To build locally

- Install `protoc`. See instruction [here](https://grpc.io/docs/protoc-installation/)
- Build:

```shell
git clone https://github.com/ozoncp/ocp-experience-api.git
cd ocp-experience-api
make build
```
Built binary exists at `bin/ocp-experience-api`

### To build and run with Docker

- Build docker image `docker build . -t ocp-experience-api`
- Run `docker run -p 82:82 ocp-experience-api`