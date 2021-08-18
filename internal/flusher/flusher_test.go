package flusher

import (
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

		journeysTable = []models.Journey{
			{JourneyId: 0, UserId: 1, Address: "Воронеж", Description: ""},
			{JourneyId: 1, UserId: 1, Address: "Уфа", Description: ""},
			{JourneyId: 2, UserId: 2, Address: "Москва", Description: ""},
			{JourneyId: 3, UserId: 2, Address: "Лондон", Description: ""},
			{JourneyId: 4, UserId: 3, Address: "Новосибирск", Description: ""},
		}
		errRepo = errors.New("repo adding error")
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(ctrl)
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

			It("should return nil and not call AddJourneys", func() {
				mockRepo.EXPECT().AddJourneys(gomock.Any()).Times(0)

				result := f.Flush(journeys)

				Expect(result).Should(BeNil())
			})
		})

		Context("journey slice is empty", func() {
			BeforeEach(func() {
				journeys = []models.Journey{}
			})

			It("should return nil and not call AddJourneys", func() {
				mockRepo.EXPECT().AddJourneys(gomock.Any()).Times(0)

				result := f.Flush(journeys)

				Expect(result).Should(BeNil())
			})
		})
	})

	Context("chunk size is negative", func() {
		BeforeEach(func() {
			chunkSize = -1
			journeys = journeysTable
		})

		It("should return original slice and not call AddJourneys", func() {
			mockRepo.EXPECT().AddJourneys(gomock.Any()).Times(0)

			result := f.Flush(journeys)

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

			It("should return nil and call AddJourneys 3 times", func() {
				mockRepo.EXPECT().AddJourneys(journeys[:2]).Times(1)
				mockRepo.EXPECT().AddJourneys(journeys[2:4]).Times(1)
				mockRepo.EXPECT().AddJourneys(journeys[4:]).Times(1)

				result := f.Flush(journeys)

				Expect(result).Should(BeNil())
			})
		})

		Context("chunkSize is greater then journeys length", func() {
			BeforeEach(func() {
				chunkSize = 7
			})

			It("should return nil and call AddJourneys 1 time", func() {
				mockRepo.EXPECT().AddJourneys(journeys).Times(1)

				result := f.Flush(journeys)

				Expect(result).Should(BeNil())
			})
		})

		Context("chunkSize is equal journeys length", func() {
			BeforeEach(func() {
				chunkSize = 5
			})

			It("should return nil and call AddJourneys 1 time", func() {
				mockRepo.EXPECT().AddJourneys(journeys).Times(1)

				result := f.Flush(journeys)

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
			It("should return second chunk and call AddJourneys 3 times", func() {
				mockRepo.EXPECT().AddJourneys(journeys[:2]).Times(1)
				mockRepo.EXPECT().AddJourneys(journeys[2:4]).Times(1).Return(errRepo)
				mockRepo.EXPECT().AddJourneys(journeys[4:]).Times(1)

				result := f.Flush(journeys)

				Expect(result).Should(Equal(journeys[2:4]))
			})
		})

		Context("on fail on every chunk", func() {
			It("should return original slice and call AddJourneys 3 times", func() {
				mockRepo.EXPECT().AddJourneys(journeys[:2]).Times(1).Return(errRepo)
				mockRepo.EXPECT().AddJourneys(journeys[2:4]).Times(1).Return(errRepo)
				mockRepo.EXPECT().AddJourneys(journeys[4:]).Times(1).Return(errRepo)

				result := f.Flush(journeys)

				Expect(result).Should(Equal(journeys))
			})
		})
	})
})
