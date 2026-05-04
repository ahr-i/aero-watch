package gps

import (
	"sync"
	"time"
)

type Position struct {
	Group     string    `json:"group"`
	Code      string    `json:"code"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store struct {
	mu        sync.RWMutex
	positions map[positionKey]Position
	ttl       time.Duration
	stopCh    chan struct{}
	doneCh    chan struct{}
}

type positionKey struct {
	group string
	code  string
}

func NewStore(ttl time.Duration, cleanupInterval time.Duration) *Store {
	store := &Store{
		positions: make(map[positionKey]Position),
		ttl:       ttl,
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
	}

	go store.cleanupLoop(cleanupInterval)

	return store
}

func (s *Store) Upsert(group string, code string, latitude float64, longitude float64) Position {
	position := Position{
		Group:     group,
		Code:      code,
		Latitude:  latitude,
		Longitude: longitude,
		UpdatedAt: time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.positions[newPositionKey(group, code)] = position

	return position
}

func (s *Store) Get(group string, code string) (Position, bool) {
	s.mu.RLock()
	key := newPositionKey(group, code)
	position, ok := s.positions[key]
	s.mu.RUnlock()
	if !ok {
		return Position{}, false
	}

	if s.isExpired(position, time.Now()) {
		s.mu.Lock()
		delete(s.positions, key)
		s.mu.Unlock()

		return Position{}, false
	}

	return position, true
}

func (s *Store) List() []Position {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	positions := make([]Position, 0, len(s.positions))
	for key, position := range s.positions {
		if s.isExpired(position, now) {
			delete(s.positions, key)
			continue
		}

		positions = append(positions, position)
	}

	return positions
}

func (s *Store) Close() {
	close(s.stopCh)
	<-s.doneCh
}

func (s *Store) cleanupLoop(cleanupInterval time.Duration) {
	defer close(s.doneCh)

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.deleteExpired()
		case <-s.stopCh:
			return
		}
	}
}

func (s *Store) deleteExpired() {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, position := range s.positions {
		if s.isExpired(position, now) {
			delete(s.positions, key)
		}
	}
}

func (s *Store) isExpired(position Position, now time.Time) bool {
	return now.Sub(position.UpdatedAt) > s.ttl
}

func newPositionKey(group string, code string) positionKey {
	return positionKey{
		group: group,
		code:  code,
	}
}
