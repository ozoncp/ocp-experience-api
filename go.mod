module github.com/ozoncp/ocp-experience-api

go 1.16

require (
	github.com/golang/mock v1.6.0 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
)

replace (
	github.com/ozoncp/ocp-request-api/pkg/ocp-experience-api => ./pkg/ocp-experience-api
)
