package api

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog/log"

	"github.com/ozoncp/ocp-experience-api/internal/models"

	repository "github.com/ozoncp/ocp-experience-api/internal/repo"
	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

// NewExperienceApi creates Experience API instance
func NewExperienceApi(r repository.Repo) *ExperienceAPI {
	return &ExperienceAPI{
		repo: r,
	}
}

type ExperienceAPI struct {
	desc.UnimplementedOcpExperienceApiServer
	repo repository.Repo
}

// ListExperienceV1 returns a list of user Requests
func (r *ExperienceAPI) ListExperienceV1(ctx context.Context, req *desc.ListExperienceV1Request) (*desc.ListExperienceV1Response, error) {
	log.Printf("ListExperienceV1 request: %v", req)

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	experiences, err := r.repo.List(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Error().Msgf("ListExperienceV1 request %v failed with %v", req, err)
		return nil, err
	}

	result := make([]*desc.Experience, 0, len(experiences))
	
	for _, experience := range experiences {
		result = append(result, models.ConvertExperienceToAPI(&experience))
	}
	
	return &desc.ListExperienceV1Response{
		Experiences: result,
	}, nil
}

// DescribeExperienceV1 returns detailed information of an experience
func (r *ExperienceAPI) DescribeExperienceV1(ctx context.Context, req *desc.DescribeExperienceV1Request) (*desc.DescribeExperienceV1Response, error) {
	log.Printf("DescribeExperienceV1 request: %v", req)

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	experience, err := r.repo.Describe(ctx, req.Id)

	if errors.Is(err, repository.NotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		log.Error().Msgf("ListExperienceV1 %v failed with %v", req, err)
		return nil, err
	}

	return &desc.DescribeExperienceV1Response{
		Experience: models.ConvertExperienceToAPI(&experience),
	}, nil
}

// CreateExperienceV1 creates new experience. Returns created object id
func (r *ExperienceAPI) CreateExperienceV1(ctx context.Context, req *desc.CreateExperienceV1Request) (*desc.CreateExperienceV1Response, error) {
	log.Printf("CreateExperienceV1 request: %v", req)

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
		log.Error().Msgf("CreateExperienceV1 %v failed with %v", req, err)
		return nil, err
	}

	return &desc.CreateExperienceV1Response{
		Id: id,
	}, nil
}

// RemoveExperienceV1 removes experience by id. Returns a removing result
func (r *ExperienceAPI) RemoveExperienceV1(ctx context.Context, req *desc.RemoveExperienceV1Request) (*desc.RemoveExperienceV1Response, error) {
	log.Printf("RemoveExperienceV1 request: %v", req)

	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := r.repo.Remove(ctx, req.Id)

	if err != nil {
		log.Error().Msgf("RemoveExperienceV1 %v failed with %v", req, err)
		return nil, err
	}

	return &desc.RemoveExperienceV1Response{
		Removed: res,
	}, nil
}
