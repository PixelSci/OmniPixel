package generation

import (
	"sync"

	"github.com/google/uuid"
)

type Generation struct {
	ConversationID uuid.UUID
	Buffer         []string
	content        string
	done           bool
	mu             sync.RWMutex
	subs           map[chan string]struct{}
}

func (g *Generation) Append(token string) {
	g.mu.Lock()
	g.Buffer = append(g.Buffer, token)
	g.content += token
	for ch := range g.subs {
		select {
		case ch <- token:
		default:
		}
	}
	g.mu.Unlock()
}

func (g *Generation) Subscribe() chan string {
	g.mu.Lock()
	defer g.mu.Unlock()
	ch := make(chan string, 256)
	g.subs[ch] = struct{}{}
	return ch
}

func (g *Generation) unsubscribe(ch chan string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.subs, ch)
}

func (g *Generation) finish() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.done = true
	for ch := range g.subs {
		close(ch)
	}
	g.subs = nil
	return g.content
}

func (g *Generation) isDone() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.done
}

func (g *Generation) Content() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.content
}

type Manager struct {
	mu    sync.RWMutex
	items map[uuid.UUID]*Generation
}

func NewManager() *Manager {
	return &Manager{items: make(map[uuid.UUID]*Generation)}
}

func (m *Manager) Start(conversationID uuid.UUID) *Generation {
	m.mu.Lock()
	defer m.mu.Unlock()
	g := &Generation{
		ConversationID: conversationID,
		subs:           make(map[chan string]struct{}),
	}
	m.items[conversationID] = g
	return g
}

func (m *Manager) Subscribe(conversationID uuid.UUID) (chan string, []string, bool) {
	m.mu.RLock()
	g, ok := m.items[conversationID]
	m.mu.RUnlock()
	if !ok || g.isDone() {
		return nil, nil, false
	}

	g.mu.RLock()
	buffer := make([]string, len(g.Buffer))
	copy(buffer, g.Buffer)
	g.mu.RUnlock()

	ch := g.Subscribe()
	return ch, buffer, true
}

func (m *Manager) Finish(conversationID uuid.UUID) (string, bool) {
	m.mu.RLock()
	g, ok := m.items[conversationID]
	m.mu.RUnlock()
	if !ok {
		return "", false
	}
	content := g.finish()
	m.mu.Lock()
	delete(m.items, conversationID)
	m.mu.Unlock()
	return content, true
}

func (m *Manager) Get(conversationID uuid.UUID) (*Generation, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	g, ok := m.items[conversationID]
	return g, ok
}

func (m *Manager) Has(conversationID uuid.UUID) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.items[conversationID]
	return ok
}
