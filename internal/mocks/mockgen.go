package mocks

//go:generate mockgen -destination=./mocks/flusher_mock.go -package=mocks github.com/ozoncp/ocp-experience-api/internal/flusher Flusher
//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks github.com/ozoncp/ocp-experience-api/internal/repo Repo
//go:generate mockgen -destination=./mocks/saver_mock.go -package=mocks github.com/ozoncp/ocp-experience-api/internal/saver Saver
//go:generate mockgen -destination=./mocks/metrics_reporter_mock.go -package=mocks github.com/ozoncp/ocp-experience-api/internal/metrics Reporter
//go:generate mockgen -destination=./mocks/producer_mock.go -package=mocks github.com/ozoncp/ocp-experience-api/internal/producer Producer
