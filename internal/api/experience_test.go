package api_test

import (
	"context"
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

var _ = Describe("Api", func() {
	var (
		experienceAPI *api.ExperienceAPI
		mockRepo      *mocks.MockRepo
		mockCtrl   *gomock.Controller
		ctx        context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(mockCtrl)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Add new item and return id", func() {
		JustBeforeEach(func() {
			experienceAPI = api.NewExperienceApi(mockRepo)
			ctx = context.Background()
		})

		It("Add request with no error", func() {
			id := uint64(11)
			mockRepo.EXPECT().
				Add(ctx, gomock.Any()).
				Return(id, nil).
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

		It("Add() params validation", func() {
			_, err := experienceAPI.CreateExperienceV1(
				ctx, &desc.CreateExperienceV1Request{UserId: 0},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid CreateExperienceV1Request.UserId: value must be greater than 0"))
		})

		It("Remove() params validation", func() {
			_, err := experienceAPI.RemoveExperienceV1(
				ctx, &desc.RemoveExperienceV1Request{},
			)

			Expect(err.Error()).To(Equal("rpc error: code = InvalidArgument desc = invalid RemoveExperienceV1Request.Id: value must be greater than 0"))
		})

		It("Describe() params validation", func() {
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
				models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
				models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
				models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
			}

			mockRepo.EXPECT().
				List(ctx, limit, offset).
				Return(requests, nil).
				Times(1)

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
			mockRepo.EXPECT().
				Remove(ctx, id).
				Return(expectFound, nil).
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

		It("Describe existing experience", func() {
			experience := models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1)

			mockRepo.EXPECT().
				Describe(ctx, experience.Id).
				Return(experience, nil).
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
				Describe(ctx, id).
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
