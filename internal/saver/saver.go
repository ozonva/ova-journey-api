package saver

import (
	"errors"
	"github.com/ozonva/ova-journey-api/internal/flusher"
	"github.com/ozonva/ova-journey-api/internal/models"
	"sync"
	"time"
)

// InternalBufferIsFullError - this error occurs when you try to add a new element using the Saver.Save() method
// if there are already capacity elements in the internal buffer
var InternalBufferIsFullError = errors.New("cannot add new journey, saver internal buffer is full")

// Saver represents the object used for saving journeys in storage
type Saver interface {
	// Save - add new models.Journey for saving
	Save(entity models.Journey) error
	// Close - save immediately all journeys from internal buffer to storage
	Close()
}

type saver struct {
	sync.Mutex
	flusher flusher.Flusher
	buffer  []models.Journey
	close   chan struct{}
}

// Save - add new journey to internal buffer of Saver
func (s *saver) Save(journey models.Journey) error {
	s.Lock()
	defer s.Unlock()

	if len(s.buffer) == cap(s.buffer) {
		return InternalBufferIsFullError
	}

	s.buffer = append(s.buffer, journey)
	return nil
}

// Close - flush all collected data from internal buffer.
func (s *saver) Close() {
	s.close <- struct{}{}
}

func (s *saver) uploadToFlusher() {
	s.Lock()
	defer s.Unlock()
	saveFailedJourneys := s.flusher.Flush(s.buffer)
	s.buffer = s.buffer[:0]
	// if not all data was flushed, restore them to try flush again in next call
	if len(saveFailedJourneys) > 0 {
		s.buffer = append(s.buffer, saveFailedJourneys...)
	}
}

// NewSaver return Saver with periodic flushing Journeys data to the storage using flusher.
//
// Use Saver.Save() method to add new journey for flushing.
// For collecting data between flushing attempts used internal buffer with capacity size.
// If internal buffer is full the Saver.Save() method returns InternalBufferIsFullError without adding journey for flushing.
//
// Use Saver.Close() method to immediately flush data from internal buffer.
func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	delayBetweenFlushing time.Duration,
) Saver {
	if capacity < 1 {
		panic("capacity must be greater then 0")
	}
	if flusher == nil {
		panic("flusher cannot be nil")
	}
	if delayBetweenFlushing < 1 {
		panic("delayBetweenFlushing must be greater then 0")
	}

	s := saver{
		flusher: flusher,
		close:   make(chan struct{}),
		buffer:  make([]models.Journey, 0, capacity),
	}

	go func() {
		ticker := time.NewTicker(delayBetweenFlushing)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.uploadToFlusher()
			case <-s.close:
				s.uploadToFlusher()
			}
		}
	}()

	return &s
}
