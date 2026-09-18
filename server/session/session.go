package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type sessionEntry struct {
	expiresAt time.Time
}

type Manager struct {
	mu             sync.Mutex
	sessions       map[string]sessionEntry
	defaultExpiry  time.Duration
	validateToken  func(string) bool
	stopCleanup    chan struct{}
}

func NewManager(defaultExpiry time.Duration) *Manager {
	m := &Manager{
		sessions:      make(map[string]sessionEntry),
		defaultExpiry: defaultExpiry,
		stopCleanup:   make(chan struct{}),
	}
	m.validateToken = func(token string) bool {
		return m.Validate(token)
	}
	go m.cleanupLoop()
	return m
}

func (m *Manager) Create() (string, time.Time) {
	return m.CreateWithExpiry(m.defaultExpiry)
}

func (m *Manager) CreateWithExpiry(expiry time.Duration) (string, time.Time) {
	token := generateToken()
	expiresAt := time.Now().Add(expiry)
	m.mu.Lock()
	m.sessions[token] = sessionEntry{expiresAt: expiresAt}
	m.mu.Unlock()
	return token, expiresAt
}

func (m *Manager) Validate(token string) bool {
	if token == "" {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	entry, ok := m.sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(entry.expiresAt) {
		delete(m.sessions, token)
		return false
	}
	return true
}

func (m *Manager) Delete(token string) {
	m.mu.Lock()
	delete(m.sessions, token)
	m.mu.Unlock()
}

func (m *Manager) ValidateFunc() func(string) bool {
	return m.validateToken
}

func (m *Manager) Stop() {
	close(m.stopCleanup)
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.stopCleanup:
			return
		}
	}
}

func (m *Manager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for token, entry := range m.sessions {
		if now.After(entry.expiresAt) {
			delete(m.sessions, token)
		}
	}
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random token: " + err.Error())
	}
	return hex.EncodeToString(b)
}
