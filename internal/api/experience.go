package api

import (
	"context"
	"math/rand"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rs/zerolog/log"

	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"
)

// NewExperienceApi creates Experience API instance
func NewExperienceApi() *ExperienceAPI {
	return &ExperienceAPI{}
}

type ExperienceAPI struct {
	desc.UnimplementedOcpExperienceApiServer
}

// ListExperienceV1 returns a list of user Requests
func (r *ExperienceAPI) ListExperienceV1(ctx context.Context, req *desc.ListExperienceV1Request) (*desc.ListExperienceV1Response, error) {
	log.Printf("ListExperienceV1 request: %v", req)

	err := req.Validate()

	if err != nil {
		return nil, err
	}

	return &desc.ListExperienceV1Response{
		Experiences: make([]*desc.Experience, 0),
	}, nil
}

// DescribeExperienceV1 returns detailed information of an experience
func (r *ExperienceAPI) DescribeExperienceV1(ctx context.Context, req *desc.DescribeExperienceV1Request) (*desc.DescribeExperienceV1Response, error) {
	err := req.Validate()

	if err != nil {
		return nil, err
	}

	log.Printf("DescribeExperienceV1 request: %v", req)

	return &desc.DescribeExperienceV1Response{
		Experience: &desc.Experience{
			Id:     1,
			UserId: 1,
			Type:   1,
			From:   &timestamp.Timestamp{},
			To:     &timestamp.Timestamp{},
			Level:  1,
		},
	}, nil
}

// CreateExperienceV1 creates new experience. Returns created object id
func (r *ExperienceAPI) CreateExperienceV1(ctx context.Context, req *desc.CreateExperienceV1Request) (*desc.CreateExperienceV1Response, error) {
	err := req.Validate()

	if err != nil {
		return nil, err
	}

	log.Printf("CreateExperienceV1 request: %v", req)

	return &desc.CreateExperienceV1Response{
		Id: rand.Uint64(),
	}, nil
}

// RemoveExperienceV1 removes experience by id. Returns a removing result
func (r *ExperienceAPI) RemoveExperienceV1(ctx context.Context, req *desc.RemoveExperienceV1Request) (*desc.RemoveExperienceV1Response, error) {
	err := req.Validate()

	if err != nil {
		return nil, err
	}

	log.Printf("RemoveExperienceV1 request: %v", req)

	return &desc.RemoveExperienceV1Response{
		Removed: false,
	}, nil
}
