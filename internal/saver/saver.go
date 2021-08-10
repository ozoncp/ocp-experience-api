package saver

import (
	"log"
	"time"

	"github.com/ozoncp/ocp-experience-api/internal/flusher"
	"github.com/ozoncp/ocp-experience-api/internal/models"
)

// saver states
const (
	saverInitialized = 0x01
	saverClosed      = 0x02
)

// Saver saves Experience into storage.
// Init() must be called before using an instance. Close() to ensure all pending item are stored.
type Saver interface {
	Save(entity models.Experience)
	Init()
	Close()
}

// NewSaver creates Saver instance. NewSaver calls flusher Init() internally.
// Async collects and save Experience entities into internally slice with given `capacity`.
// duration represents tick range
func NewSaver(capacity uint, flusher flusher.Flusher, duration time.Duration) Saver {
	s := &saver{
		capacity:     capacity,
		flusher:      flusher,
		queueChan:    make(chan models.Experience, capacity),
		entities:     make([]models.Experience, 0, capacity),
		tickDuration: duration,
		closeChan:    make(chan struct{}),
	}

	return s
}

// Implements Saver interface
type saver struct {
	capacity     uint
	flusher      flusher.Flusher
	queueChan    chan models.Experience
	entities     []models.Experience
	state        int8 // state may be initialized or closed
	tickDuration time.Duration
	closeChan    chan struct{}
}

// Save saves Experience into storage
func (s *saver) Save(entity models.Experience) {
	s.assertNotClosed()
	s.assertInitialized()

	s.queueChan <- entity
}

// Init inits saver, this method should be called before Saver usage
func (s *saver) Init() {
	s.assertNotClosed()

	if s.isInitialized() {
		return
	}

	go s.run()

	s.setState(saverInitialized)
}

// Close closes saver. Ensures that all entity object are processed.
func (s *saver) Close() {
	s.assertInitialized()

	if s.isClosed() {
		return
	}

	s.closeChan <- struct{}{}
}

// main loop
func (s *saver) run() {
	timer := time.NewTicker(s.tickDuration)
	defer timer.Stop()

	for {
		select {
		case res := <-s.queueChan:
			s.entities = append(s.entities, res)

		case <-timer.C:
			s.flush()

		case <-s.closeChan:
			close(s.queueChan)
			close(s.closeChan)
			s.setState(saverClosed)

			return
		}
	}
}

// returns true if saver closed
func (s *saver) isClosed() bool {
	return (s.state & saverClosed) == saverClosed
}

// returns true if saver initialized
func (s *saver) isInitialized() bool {
	return (s.state & saverInitialized) == saverInitialized
}

// sets state by |= flag
func (s *saver) setState(state int8) {
	s.state |= state
}

// flush flushes Experiences to flusher
func (s *saver) flush() {
	remains, err := s.flusher.Flush(s.entities)

	if err != nil {
		log.Printf("Failed to save %v experience entities: %v", len(remains), err)
		return
	}

	s.entities = remains
}

// asserts on closed state
func (s *saver) assertNotClosed() {
	if s.isClosed() {
		panic("Saver is closed")
	}
}

// asserts on init state
func (s *saver) assertInitialized() {
	if !s.isInitialized() {
		panic("Saver is not initialized")
	}
}
