package flusher_test

import (
	"context"
	"errors"
	"time"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-experience-api/internal/flusher"
	"github.com/ozoncp/ocp-experience-api/internal/mocks/mocks"
	"github.com/ozoncp/ocp-experience-api/internal/models"
)

var _ = Describe("Flusher", func() {
	var (
		flusherImpl flusher.Flusher
		mockRepo    *mocks.MockRepo
		mockCtrl    *gomock.Controller
		ctx			context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Add items with no errors", func() {
		JustBeforeEach(func() {
			flusherImpl = flusher.NewFlusher(2, mockRepo)
		})

		It("Added all bulks, one call to repo. No remains left.", func() {
			mockRepo.EXPECT().
				AddExperiences(ctx, gomock.Any()).
				Return(nil).
				MaxTimes(1).
				MinTimes(1)

			remains, err := flusherImpl.Flush(ctx, []models.Experience{
				models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
				models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
			})

			Expect(remains).To(HaveLen(0))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Added all bulks, two calls against repo. No remains left.", func() {
			mockRepo.EXPECT().
				AddExperiences(ctx, gomock.Any()).
				Return(nil).
				MaxTimes(2).
				MinTimes(2)

			remains, err := flusherImpl.Flush(ctx, []models.Experience{
				models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
				models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
				models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
				models.NewExperience(4, 4, 4, time.Time{}, time.Time{}, 4),
			})

			Expect(remains).To(HaveLen(0))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Added all bulks, two calls against repo. 1 item remained.", func() {
			mockRepo.EXPECT().
				AddExperiences(ctx, gomock.Any()).
				Return(nil).
				MaxTimes(2).
				MinTimes(2)

			remains, err := flusherImpl.Flush(ctx, []models.Experience{
				models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
				models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
				models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
				models.NewExperience(4, 4, 4, time.Time{}, time.Time{}, 4),
				models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5),
			})

			Expect(remains).To(Equal([]models.Experience{models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5)}))
			Expect(err).ToNot(HaveOccurred())
		})

	})

	Context("Failed to add items to Repo", func() {
		JustBeforeEach(func() {
			flusherImpl = flusher.NewFlusher(2, mockRepo)
		})

		It("Failed to add all items", func() {
			mockRepo.EXPECT().
				AddExperiences(ctx, gomock.Any()).
				Return(errors.New("failed to add")).
				MaxTimes(1).
				MinTimes(1)

			requests := []models.Experience{
				models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
				models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
			}

			remains, err := flusherImpl.Flush(ctx, requests)

			Expect(remains).To(Equal(requests))
			Expect(err).To(HaveOccurred())
		})

		It("Partially failed to add items", func() {
			successFullCall1 := mockRepo.EXPECT().
				AddExperiences(ctx, gomock.Any()).
				Return(nil)

			successFullCall2 := mockRepo.EXPECT().
				AddExperiences(ctx, gomock.Any()).
				Return(nil)

			failedCall := mockRepo.EXPECT().
				AddExperiences(ctx, gomock.Any()).
				Return(errors.New("failed to add"))

			gomock.InOrder(successFullCall1, successFullCall2, failedCall)

			requests := []models.Experience{
				models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
				models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
				models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
				models.NewExperience(4, 4, 4, time.Time{}, time.Time{}, 4),
				models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5),
				models.NewExperience(6, 6, 6, time.Time{}, time.Time{}, 6),
				models.NewExperience(7, 7, 7, time.Time{}, time.Time{}, 7),
			}

			remains, err := flusherImpl.Flush(ctx, requests)

			Expect(remains).To(Equal([]models.Experience{
				models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5),
				models.NewExperience(6, 6, 6, time.Time{}, time.Time{}, 6),
				models.NewExperience(7, 7, 7, time.Time{}, time.Time{}, 7),
			}), "These are failed to add to repo")

			Expect(err).To(HaveOccurred())
		})
	})
})
