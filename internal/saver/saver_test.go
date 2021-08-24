package saver_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-experience-api/internal/mocks/mocks"
	"github.com/ozoncp/ocp-experience-api/internal/models"
	"github.com/ozoncp/ocp-experience-api/internal/saver"
)

// creates test data
func makeExperienceEntities(num uint64) []models.Experience {
	entities := make([]models.Experience, 0, num)

	for i := uint64(0); i < num; i++ {
		entities = append(entities, models.NewExperience(i, i, i, time.Time{}, time.Time{}, i))
	}

	return entities
}

var _ = Describe("Saver", func() {
	var (
		s           saver.Saver
		mockFlusher *mocks.MockFlusher
		mockCtrl    *gomock.Controller
		entities    []models.Experience
		ctx         context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockCtrl = gomock.NewController(GinkgoT())
		mockFlusher = mocks.NewMockFlusher(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Saver test", func() {
		JustBeforeEach(func() {
			entities = makeExperienceEntities(10)
		})

		It("All items are saved after tick", func() {
			s = saver.NewSaver(10, mockFlusher, time.Millisecond * 200)
			s.Init()
			defer s.Close()

			mockFlusher.EXPECT().Flush(ctx, gomock.Eq(entities)).Times(1).Return(nil, nil)

			for _, e := range entities {
				s.Save(e)
			}

			time.Sleep(time.Millisecond * 250)
		})

		It("All items are saved after Saver flushes with two intervals", func() {
			s = saver.NewSaver(10, mockFlusher, time.Millisecond * 200)
			s.Init()
			defer s.Close()

			callFirst := mockFlusher.EXPECT().Flush(ctx, gomock.Eq(entities[:len(entities)/2])).Times(1).Return(nil, nil)
			mockFlusher.EXPECT().Flush(ctx, gomock.Eq(entities[len(entities)/2:])).Times(1).Return(nil, nil).After(callFirst)

			for _, e := range entities[:len(entities)/2] {
				s.Save(e)
			}

			time.Sleep(time.Millisecond * 250)

			for _, e := range entities[len(entities)/2:] {
				s.Save(e)
			}

			time.Sleep(time.Millisecond * 200)
		})

		It("Items has not been saved, closed earlier", func() {
			s = saver.NewSaver(10, mockFlusher, time.Millisecond * 500)
			s.Init()
			defer s.Close()

			mockFlusher.EXPECT().Flush(ctx, gomock.Any()).Times(0)

			for _, e := range entities {
				s.Save(e)
			}
		})
	})

	Context("Saver state assertions test", func() {
		JustBeforeEach(func() {
			s = saver.NewSaver(10, mockFlusher, time.Second)
		})

		It("Must call Init() before", func() {
			Expect(func() {
				s.Save(models.Experience{})
			}).To(PanicWith("Saver is not initialized"))
		})

		It("Cannot Save() after Close()", func() {
			s.Init()
			s.Close()

			Expect(func() {
				s.Save(models.Experience{})
			}).To(PanicWith("Saver is closed"))
		})

		It("Cannot Init() after Close()", func() {
			s.Init()
			s.Close()

			Expect(func() {
				s.Init()
			}).To(PanicWith("Saver is closed"))
		})

		It("Cannot Close() before Init()", func() {
			Expect(func() {
				s.Close()
			}).To(PanicWith("Saver is not initialized"))
		})
	})
})
