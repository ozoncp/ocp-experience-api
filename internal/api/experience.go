package api

import (
	"context"
	"errors"
	"github.com/ozoncp/ocp-experience-api/internal/utils"

	"github.com/opentracing/opentracing-go"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog/log"

	"github.com/ozoncp/ocp-experience-api/internal/metrics"
	"github.com/ozoncp/ocp-experience-api/internal/models"
	"github.com/ozoncp/ocp-experience-api/internal/producer"

	traceLog "github.com/opentracing/opentracing-go/log"
	repository "github.com/ozoncp/ocp-experience-api/internal/repo"
	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

type validator interface {
	Validate() error
}

// NewExperienceApi creates Experience API instance
func NewExperienceApi(r repository.Repo,
	batchSize uint64,
	reporter metrics.Reporter,
	producer producer.Producer,
	tracer opentracing.Tracer) *ExperienceAPI {

	return &ExperienceAPI{
		repo: r,
		batchSize: batchSize,
		metrics : reporter,
		producer: producer,
		tracer: tracer,
	}
}

type ExperienceAPI struct {
	desc.UnimplementedOcpExperienceApiServer
	repo      repository.Repo
	batchSize uint64
	metrics   metrics.Reporter
	producer  producer.Producer
	tracer    opentracing.Tracer
}

// ListExperienceV1 returns a list of user Requests
func (r *ExperienceAPI) ListExperienceV1(ctx context.Context, req *desc.ListExperienceV1Request) (*desc.ListExperienceV1Response, error) {
	log.Printf("ListExperienceV1 request: %v", req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "ListExperienceV1")
	defer span.Finish()

	if err := r.validate(ctx, req, producer.ReadEvent); err != nil {
		return nil, err
	}

	experiences, err := r.repo.List(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Error().
			Err(err).
			Str("endpoint", "ListExperienceV1").
			Uint64("limit", req.Limit).
			Uint64("offset", req.Offset).
			Msgf("Failed to list experiences")

		r.producer.Send(producer.NewEvent(ctx, 0, producer.ReadEvent, err))
		return nil, err
	}

	result := make([]*desc.Experience, 0, len(experiences))
	eventMessages := make([]producer.EventMsg, 0, len(experiences))

	for _, experience := range experiences {
		result = append(result, models.ConvertExperienceToAPI(&experience))
		eventMessages = append(eventMessages, producer.NewEvent(ctx, experience.Id, producer.ReadEvent, nil))

		r.producer.Send(eventMessages...)
	}

	r.metrics.IncList(1, "ListExperienceV1")
	return &desc.ListExperienceV1Response{
		Experiences: result,
	}, nil
}

// DescribeExperienceV1 returns detailed information of an experience
func (r *ExperienceAPI) DescribeExperienceV1(ctx context.Context, req *desc.DescribeExperienceV1Request) (*desc.DescribeExperienceV1Response, error) {
	log.Printf("DescribeExperienceV1 request: %v", req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "DescribeExperienceV1")
	defer span.Finish()

	if err := r.validate(ctx, req, producer.ReadEvent); err != nil {
		return nil, err
	}

	experience, err := r.repo.Describe(ctx, req.Id)

	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		log.Error().
			Str("endpoint", "DescribeExperienceV1").
			Uint64("id", req.Id).
			Err(err).
			Msgf("Failed to read experience")

		return nil, err
	}

	r.producer.Send(producer.NewEvent(ctx, req.Id, producer.ReadEvent, err))
	r.metrics.IncRead(1, "DescribeExperienceV1")

	return &desc.DescribeExperienceV1Response{
		Experience: models.ConvertExperienceToAPI(&experience),
	}, nil
}

// CreateExperienceV1 creates new experience. Returns created object id
func (r *ExperienceAPI) CreateExperienceV1(ctx context.Context, req *desc.CreateExperienceV1Request) (*desc.CreateExperienceV1Response, error) {
	log.Printf("CreateExperienceV1 request: %v", req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateExperienceV1")
	defer span.Finish()

	if err := r.validate(ctx, req, producer.CreateEvent); err != nil {
		return nil, err
	}

	id, err := r.repo.Add(ctx, models.Experience{
		Id:     0,
		UserId: req.UserId,
		Type:   req.Type,
		From:   req.From.AsTime(),
		To:     req.To.AsTime(),
		Level:  req.Level,
	})

	if err != nil {
		log.Error().
			Str("endpoint", "CreateExperienceV1").
			Err(err).
			Msgf("Failed to create experience")

		return nil, err
	}

	r.producer.Send(producer.NewEvent(ctx, id, producer.CreateEvent, err))
	r.metrics.IncCreate(1, "CreateExperienceV1")

	return &desc.CreateExperienceV1Response{
		Id: id,
	}, nil
}

// MultiCreateExperienceV1  Creates new experiences, returns new ids
func (r *ExperienceAPI) MultiCreateExperienceV1(ctx context.Context, req *desc.MultiCreateExperienceV1Request) (*desc.MultiCreateExperienceV1Response, error) {
	log.Printf("Multi create experience: %v", req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "MultiCreateExperienceV1")
	defer span.Finish()

	if err := r.validate(ctx, req, producer.CreateEvent); err != nil {
		return nil, err
	}

	toCreate := make([]models.Experience, 0, len(req.Experiences))

	for _, experience := range req.Experiences {
		toCreate = append(toCreate, models.NewExperience(0, experience.UserId, experience.Type, experience.From.AsTime(), experience.To.AsTime(), experience.Level))
	}

	newIds := make([]uint64, 0, len(req.Experiences))
	batch, err := utils.SplitExperienceToBulks(toCreate, int(r.batchSize))

	if err != nil {
		return nil, err
	}

	for _, batch := range batch {
		ids, writeErr := r.writeExperiencesBatch(ctx, batch)

		if writeErr != nil {
			return nil, writeErr
		}

		newIds = append(newIds, ids...)
		r.metrics.IncCreate(uint(len(ids)), "MultiCreateExperienceV1")
	}

	return &desc.MultiCreateExperienceV1Response{
		Ids: newIds,
	}, nil
}

// RemoveExperienceV1 removes experience by id. Returns a removing result
func (r *ExperienceAPI) RemoveExperienceV1(ctx context.Context, req *desc.RemoveExperienceV1Request) (*desc.RemoveExperienceV1Response, error) {
	log.Printf("RemoveExperienceV1 request: %v", req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "RemoveExperienceV1")
	defer span.Finish()

	if err := r.validate(ctx, req, producer.DeleteEvent); err != nil {
		return nil, err
	}

	res, err := r.repo.Remove(ctx, req.Id)

	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, "experience does not exist")
	}

	if err != nil {
		log.Error().
			Err(err).
			Uint64("id", req.Id).
			Str("endpoint", "RemoveExperienceV1").
			Msgf("Failed to remove experience")

		return nil, err
	}

	r.producer.Send(producer.NewEvent(ctx, req.Id, producer.DeleteEvent, err))
	r.metrics.IncRemove(1, "RemoveExperienceV1")

	return &desc.RemoveExperienceV1Response{
		Removed: res,
	}, nil
}

// UpdateExperienceV1 updates experience
func (r *ExperienceAPI) UpdateExperienceV1(ctx context.Context, req *desc.UpdateExperienceV1Request) (*desc.UpdateExperienceV1Response, error) {
	log.Printf("Update request: %v", req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateExperienceV1")
	defer span.Finish()

	if err := r.validate(ctx, req, producer.UpdateEvent); err != nil {
		return nil, err
	}

	err := r.repo.Update(ctx, models.NewExperience(req.Id, req.UserId, req.Type, req.From.AsTime(), req.To.AsTime(), req.Level))

	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, "experience does not exist")
	}

	if err != nil {
		log.Error().
			Uint64("id", req.Id).
			Str("endpoint", "UpdateExperienceV1").
			Err(err).
			Msgf("Failed to update experience")

		return nil, err
	}

	r.producer.Send(producer.NewEvent(ctx, req.Id, producer.UpdateEvent, err))
	r.metrics.IncUpdate(1, "UpdateExperienceV1")

	return &desc.UpdateExperienceV1Response{}, nil
}

func (r *ExperienceAPI) validate(ctx context.Context, request validator, event producer.EventType) error {
	if err := request.Validate(); err != nil {
		r.producer.Send(producer.NewEvent(ctx, 0, event, err))
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func (r *ExperienceAPI) writeExperiencesBatch(ctx context.Context, batch []models.Experience) ([]uint64, error) {
	childSpan, childCtx := opentracing.StartSpanFromContext(ctx, "MultiCreateExperienceV1Batch")
	childSpan.LogFields(traceLog.Int("batch_size", len(batch)))
	defer childSpan.Finish()

	ids, err := r.repo.AddExperiences(childCtx, batch)

	if err != nil {
		log.Error().Err(err).Msgf("Failed to save experiences")
		r.producer.Send(producer.NewEvent(ctx, 0, producer.CreateEvent, err))

		return nil, err
	}

	for _, id := range ids {
		r.producer.Send(producer.NewEvent(ctx, id, producer.CreateEvent, nil))
	}

	return ids, nil
}
