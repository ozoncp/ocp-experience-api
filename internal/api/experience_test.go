package api_test

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-experience-api/internal/api"
	"github.com/ozoncp/ocp-experience-api/internal/mocks/mocks"
	"github.com/ozoncp/ocp-experience-api/internal/models"
	"github.com/ozoncp/ocp-experience-api/internal/repo"

	desc "github.com/ozoncp/ocp-experience-api/pkg/ocp-experience-api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AnyContextType struct {}

func (c *AnyContextType) Matches(val interface{}) bool {
	_, ok := val.(context.Context)
	return ok
}

func (c *AnyContextType) String() string {
	return "Asserts parameter is a context.Context type"
}

var _ = Describe("Api", func() {
	var (
		experienceAPI 	*api.ExperienceAPI
		mockRepo      	*mocks.MockRepo
		mockCtrl   		*gomock.Controller
		ctx        		context.Context
		mockProm     	*mocks.MockReporter
		mockProducer 	*mocks.MockProducer
	)

	var ctxType AnyContextType

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(mockCtrl)
		mockProm = mocks.NewMockReporter(mockCtrl)
		mockProducer = mocks.NewMockProducer(mockCtrl)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Add new item and return id", func() {
		JustBeforeEach(func() {
			experienceAPI = api.NewExperienceApi(
				mockRepo,
				2,
				mockProm,
				mockProducer,
				opentracing.NoopTracer{},
			)
			ctx = context.Background()
		})

		It("Add with no error", func() {
			id := uint64(11)
			mockRepo.EXPECT().
				Add(ctxType, gomock.Any()).
				Return(id, nil).
				Times(1)

			mockProm.EXPECT().
				IncCreate(uint(1), "CreateExperienceV1").
				Times(1)

			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(1)

			resp, err := experienceAPI.CreateExperienceV1(
				ctx, &desc.CreateExperienceV1Request{
					UserId: 1,
					Type:   1,
					From:   nil,
					To:     nil,
					Level:  1,
				},
			)

			Expect(resp).
				To(Equal(&desc.CreateExperienceV1Response{
					Id: id,
				}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("Add slice experience with no error", func() {
			experiences := []models.Experience{
				models.NewExperience(0, 1, 1, time.Now(), time.Now(), 1),
				models.NewExperience(0, 2, 2, time.Now(), time.Now(), 2),
				models.NewExperience(0, 3, 3, time.Now(), time.Now(), 3),
			}

			createExperienceV1Requests := make([]*desc.CreateExperienceV1Request, 0)

			for _, r := range experiences {
				createExperienceV1Requests = append(createExperienceV1Requests, &desc.CreateExperienceV1Request{
					UserId: r.UserId,
					Type:   r.Type,
					From: 	timestamppb.New(r.From),
					To: 	timestamppb.New(r.To),
					Level: r.Level,
				})
			}

			mockRepo.EXPECT().
				AddExperiences(ctxType, experiences[:2]).
				Return([]uint64{1, 2}, nil).
				Times(1)

			mockRepo.EXPECT().
				AddExperiences(ctxType, experiences[2:]).
				Return([]uint64{3}, nil).
				Times(1)

			mockProm.EXPECT().
				IncCreate(uint(2), "MultiCreateExperienceV1").
				Times(1)

			mockProm.EXPECT().
				IncCreate(uint(1), "MultiCreateExperienceV1").
				Times(1)

			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(3)

			resp, err := experienceAPI.MultiCreateExperienceV1(
				ctx, &desc.MultiCreateExperienceV1Request {
					Experiences: createExperienceV1Requests,
				},
			)

			Expect(resp).
				To(Equal(&desc.MultiCreateExperienceV1Response{
					Ids: []uint64{1, 2, 3},
				}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("Add() params validation", func() {
			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(1)

			_, err := experienceAPI.CreateExperienceV1(
				ctx, &desc.CreateExperienceV1Request{UserId: 0},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid CreateExperienceV1Request.UserId: value must be greater than 0"))
		})

		It("Remove() params validation", func() {
			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(1)

			_, err := experienceAPI.RemoveExperienceV1(
				ctx, &desc.RemoveExperienceV1Request{},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid RemoveExperienceV1Request.Id: value must be greater than 0"))
		})

		It("Describe() params validation", func() {
			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(1)

			_, err := experienceAPI.DescribeExperienceV1(
				ctx, &desc.DescribeExperienceV1Request{},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid DescribeExperienceV1Request.Id: value must be greater than 0"))
		})

		It("List() params validation", func() {
			_, err := experienceAPI.ListExperienceV1(
				ctx, &desc.ListExperienceV1Request{Limit: 0},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid ListExperienceV1Request.Limit: value must be inside range (0, 10000]"))

			_, err = experienceAPI.ListExperienceV1(
				ctx, &desc.ListExperienceV1Request{Limit: 1000000},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid ListExperienceV1Request.Limit: value must be inside range (0, 10000]"))
		})

		It("List experience with no error", func() {
			offset, limit := uint64(10), uint64(100)
			requests := []models.Experience{
				models.NewExperience(1, 1, 1, time.Now(), time.Now(), 1),
				models.NewExperience(2, 2, 2, time.Now(), time.Now(), 2),
				models.NewExperience(3, 3, 3, time.Now(), time.Now(), 3),
			}

			mockProm.EXPECT().
				IncList(uint(1), "ListExperienceV1").
				Times(1)

			mockRepo.EXPECT().
				List(ctxType, limit, offset).
				Return(requests, nil).
				Times(1)

			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(3)

			resp, err := experienceAPI.ListExperienceV1(
				ctx, &desc.ListExperienceV1Request{
					Offset: offset, Limit: limit,
				},
			)

			req := make([]*desc.Experience, 0, len(requests))

			for _, r := range requests {
				req = append(req, models.ConvertExperienceToAPI(&r))
			}

			Expect(resp).
				To(Equal(&desc.ListExperienceV1Response{
					Experiences: req,
				}))

			Expect(err).ToNot(HaveOccurred())
		})

		removeTest := func(expectFound bool) {
			id := uint64(11)

			if expectFound {
				mockRepo.EXPECT().
					Remove(ctxType, id).
					Return(nil).
					Times(1)
			} else {
				mockRepo.EXPECT().
					Remove(ctxType, id).
					Return(repo.NotFound, nil).
					Times(1)
			}

			mockProm.EXPECT().
				IncRemove(uint(1), "RemoveExperienceV1").
				Times(1)

			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(1)

			resp, err := experienceAPI.RemoveExperienceV1(
				ctx, &desc.RemoveExperienceV1Request{
					Id: id,
				},
			)

			Expect(resp).
				To(Equal(&desc.RemoveExperienceV1Response{
					Removed: expectFound,
				}))

			Expect(err).ToNot(HaveOccurred())
		}


		It("Remove existing experience with no errors", func() {
			removeTest(true)
		})

		It("Remove non-existing experience with no errors", func() {
			removeTest(false)
		})

		It("Update existing experience", func() {
			req := models.NewExperience(1, 1, 1, time.Now(), time.Now(), 1)
			mockRepo.EXPECT().
				Update(ctxType, req).
				Return(nil).
				Times(1)

			mockProm.EXPECT().
				IncUpdate(uint(1), "UpdateExperienceV1").
				Times(1)

			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(1)

			resp, err := experienceAPI.UpdateExperienceV1(
				ctx, &desc.UpdateExperienceV1Request{
					Id: req.Id,
					UserId:    	req.UserId,
					Type:      	req.Type,
					From:      	timestamppb.New(req.From),
					To: 		timestamppb.New(req.To),
					Level: 		req.Level,
				},
			)

			Expect(resp).To(Equal(&desc.UpdateExperienceV1Response{}))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Describe existing experience", func() {
			experience := models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1)

			mockRepo.EXPECT().
				Describe(ctxType, experience.Id).
				Return(experience, nil).
				Times(1)

			mockProm.EXPECT().
				IncRead(uint(1), "DescribeExperienceV1").
				Times(1)

			mockProducer.EXPECT().
				Send(gomock.Any()).
				Times(1)

			resp, err := experienceAPI.DescribeExperienceV1(
				ctx, &desc.DescribeExperienceV1Request{
					Id: experience.Id,
				},
			)

			Expect(resp).
				To(Equal(&desc.DescribeExperienceV1Response{
					Experience: models.ConvertExperienceToAPI(&experience),
				}))

			Expect(err).ToNot(HaveOccurred())
		})

		It("Describe no existing experience", func() {
			id := uint64(11)
			mockRepo.EXPECT().
				Describe(ctxType, id).
				Return(models.Experience{}, repo.NotFound).
				Times(1)

			resp, err := experienceAPI.DescribeExperienceV1(
				ctx, &desc.DescribeExperienceV1Request{
					Id: id,
				},
			)

			Expect(resp).To(BeNil())
			Expect(err).To(Equal(status.Error(codes.NotFound, repo.NotFound.Error())))
		})
	})
})
