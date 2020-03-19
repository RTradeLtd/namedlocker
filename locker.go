// Package namedlocker implements in-memory named locks.
package namedlocker

import (
	"errors"
	"sync"
)

var (
	// ErrUnlockOfUnlockedKey is the error reported when unlocking an unlocked key.
	ErrUnlockOfUnlockedKey = errors.New("unlock of unlocked key")
)

type ref struct {
	sync.RWMutex
}

// Store is an in-memory store of named locks.
//
// The zero-value is not usable
type Store struct {
	mu   sync.RWMutex
	refs map[string]*ref
}

// New creates a new named locker
func New() *Store {
	return &Store{refs: make(map[string]*ref)}
}

// Lock acquires a lock on key.
// If key is locked, it blocks until it can be acquired.
func (s *Store) Lock(key string) {
	s.mu.Lock()
	if _, ok := s.refs[key]; !ok {
		s.refs[key] = new(ref)
	}
	s.mu.Unlock()
	s.refs[key].Lock()
}

// RLock acquires a read lock on key.
// If key is locked, it blocks until it can be acquired.
func (s *Store) RLock(key string) {
	s.mu.Lock()
	if _, ok := s.refs[key]; !ok {
		s.refs[key] = new(ref)
	}
	s.mu.Unlock()
	s.refs[key].RLock()
}

// RUnlock is a wrapper around TryRUnlock that panics if it returns an error
func (s *Store) RUnlock(key string) {
	if err := s.TryRUnlock(key); err != nil {
		panic(err)
	}
}

// Unlock is a wrapper around TryUnlock that panics if it returns an error.
func (s *Store) Unlock(key string) {
	if err := s.TryUnlock(key); err != nil {
		panic(err)
	}
}

// TryUnlock releases the lock on key.
//
// If key is not locked, ErrUnlockOfUnlockedKey is returned.
func (s *Store) TryUnlock(key string) error {
	s.mu.Lock()
	if _, ok := s.refs[key]; !ok {
		s.mu.Unlock()
		return ErrUnlockOfUnlockedKey
	}
	s.refs[key].Unlock()
	s.mu.Unlock()
	return nil
}

// TryRUnlock releases the read lock on key.
//
// If key is not locked, ErrUnlockOfUnlockedKey is returned.
func (s *Store) TryRUnlock(key string) error {
	s.mu.Lock()
	if _, ok := s.refs[key]; !ok {
		s.mu.Unlock()
		return ErrUnlockOfUnlockedKey
	}
	s.refs[key].RUnlock()
	s.mu.Unlock()
	return nil
}
