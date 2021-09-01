package flusher

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozonva/ova-journey-api/internal/mocks"
	"github.com/ozonva/ova-journey-api/internal/models"
)

var _ = Describe("Flusher", func() {
	var (
		ctrl      *gomock.Controller
		mockRepo  *mocks.MockRepo
		f         Flusher
		journeys  []models.Journey
		chunkSize int
		ctx       context.Context

		journeysTable = []models.Journey{
			{JourneyID: 0, UserID: 1, Address: "Воронеж", Description: ""},
			{JourneyID: 1, UserID: 1, Address: "Уфа", Description: ""},
			{JourneyID: 2, UserID: 2, Address: "Москва", Description: ""},
			{JourneyID: 3, UserID: 2, Address: "Лондон", Description: ""},
			{JourneyID: 4, UserID: 3, Address: "Новосибирск", Description: ""},
		}
		errRepo = errors.New("repo adding error")
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(ctrl)
		ctx = context.Background()
	})

	JustBeforeEach(func() {
		f = NewFlusher(chunkSize, mockRepo)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("journey slice is nil or empty", func() {
		BeforeEach(func() {
			chunkSize = 1
		})

		Context("journey slice is nil", func() {
			BeforeEach(func() {
				journeys = nil
			})

			It("should return nil and not call AddJourneysMulti", func() {
				mockRepo.EXPECT().AddJourneysMulti(ctx, gomock.Any()).Times(0)

				result := f.Flush(ctx, journeys)

				Expect(result).Should(BeNil())
			})
		})

		Context("journey slice is empty", func() {
			BeforeEach(func() {
				journeys = []models.Journey{}
			})

			It("should return nil and not call AddJourneysMulti", func() {
				mockRepo.EXPECT().AddJourneysMulti(ctx, gomock.Any()).Times(0)

				result := f.Flush(ctx, journeys)

				Expect(result).Should(BeNil())
			})
		})
	})

	Context("chunk size is negative", func() {
		BeforeEach(func() {
			chunkSize = -1
			journeys = journeysTable
		})

		It("should return original slice and not call AddJourneysMulti", func() {
			mockRepo.EXPECT().AddJourneysMulti(ctx, gomock.Any()).Times(0)

			result := f.Flush(ctx, journeys)

			Expect(result).Should(Equal(journeys))
		})
	})

	Context("valid chunkSize", func() {
		BeforeEach(func() {
			journeys = journeysTable
		})

		Context("chunkSize is less then journeys length", func() {
			BeforeEach(func() {
				chunkSize = 2
			})

			It("should return nil and call AddJourneysMulti 3 times", func() {
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[:2]).Times(1)
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[2:4]).Times(1)
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[4:]).Times(1)

				result := f.Flush(ctx, journeys)

				Expect(result).Should(BeNil())
			})
		})

		Context("chunkSize is greater then journeys length", func() {
			BeforeEach(func() {
				chunkSize = 7
			})

			It("should return nil and call AddJourneysMulti 1 time", func() {
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys).Times(1)

				result := f.Flush(ctx, journeys)

				Expect(result).Should(BeNil())
			})
		})

		Context("chunkSize is equal journeys length", func() {
			BeforeEach(func() {
				chunkSize = 5
			})

			It("should return nil and call AddJourneysMulti 1 time", func() {
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys).Times(1)

				result := f.Flush(ctx, journeys)

				Expect(result).Should(BeNil())
			})
		})

	})

	Context("errors in repo", func() {
		BeforeEach(func() {
			chunkSize = 2
			journeys = journeysTable
		})

		Context("on fail on second chunk", func() {
			It("should return second chunk and call AddJourneysMulti 3 times", func() {
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[:2]).Times(1)
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[2:4]).Times(1).Return(errRepo)
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[4:]).Times(1)

				result := f.Flush(ctx, journeys)

				Expect(result).Should(Equal(journeys[2:4]))
			})
		})

		Context("on fail on every chunk", func() {
			It("should return original slice and call AddJourneysMulti 3 times", func() {
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[:2]).Times(1).Return(errRepo)
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[2:4]).Times(1).Return(errRepo)
				mockRepo.EXPECT().AddJourneysMulti(ctx, journeys[4:]).Times(1).Return(errRepo)

				result := f.Flush(ctx, journeys)

				Expect(result).Should(Equal(journeys))
			})
		})
	})
})
