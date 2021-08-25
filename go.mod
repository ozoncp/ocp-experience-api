module github.com/ozoncp/ocp-experience-api

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/Masterminds/squirrel v1.5.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
	github.com/tkanos/gonfig v0.0.0-20210106201359-53e13348de2f // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210813162853-db860fec028c
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/ozoncp/ocp-request-api/pkg/ocp-experience-api => ./pkg/ocp-experience-api
