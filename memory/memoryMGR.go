package memory

import (
	uuid "github.com/satori/go.uuid"
	"github.com/ilyzz/session/Structure"
	"github.com/ilyzz/session/errors"
	"sync"
)

type MemorySessionMGR struct {
	SessionMap map[string]Structure.Session
	rwlock     sync.RWMutex
}

func NewMemorySession(id string) *MemorySession {
	return &MemorySession{
		data:   make(map[string]interface{},16),
		id:     id,
		rwlock: sync.RWMutex{},
	}
}

func (m *MemorySessionMGR) Get(sessionId string) (session Structure.Session, err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	session, ok := m.SessionMap[sessionId]
	if !ok {
		err = errors.ErrSessionNotExist
		return
	}
	return
}

func (m *MemorySessionMGR) CreateSession() (session Structure.Session, err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	id := uuid.NewV4()
	sessionID := id.String()

	session = NewMemorySession(sessionID)

	m.SessionMap[sessionID] = session
	return
}


