package codexoauth

import (
	"errors"
	"sync"
	"time"

	"github.com/songquanpeng/one-api/model"
)

const refreshBuffer = time.Minute

type cachedAccessToken struct {
	token     string
	expiresAt time.Time
}

func (token cachedAccessToken) valid() bool {
	return token.token != "" && time.Until(token.expiresAt) > refreshBuffer
}

type Manager struct {
	mu           sync.RWMutex
	accessTokens map[string]cachedAccessToken
	refreshLocks map[string]*sync.Mutex
}

var DefaultManager = NewManager()

func NewManager() *Manager {
	return &Manager{
		accessTokens: make(map[string]cachedAccessToken),
		refreshLocks: make(map[string]*sync.Mutex),
	}
}

func (m *Manager) GetValidTokenForAccount(accountID string) (string, error) {
	if accountID == "" {
		return "", errors.New("codex oauth account_id is required")
	}

	m.mu.RLock()
	cached, ok := m.accessTokens[accountID]
	m.mu.RUnlock()
	if ok && cached.valid() {
		return cached.token, nil
	}

	lock := m.refreshLock(accountID)
	lock.Lock()
	defer lock.Unlock()

	m.mu.RLock()
	cached, ok = m.accessTokens[accountID]
	m.mu.RUnlock()
	if ok && cached.valid() {
		return cached.token, nil
	}

	account, err := model.GetCodexOAuthAccountByAccountID(accountID)
	if err != nil {
		return "", err
	}
	return m.refreshAccount(account)
}

func (m *Manager) ClearAccountToken(accountID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.accessTokens, accountID)
}

func (m *Manager) StoreAccessToken(accountID string, token string, expiresAt time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.accessTokens[accountID] = cachedAccessToken{token: token, expiresAt: expiresAt}
}

func (m *Manager) refreshLock(accountID string) *sync.Mutex {
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, ok := m.refreshLocks[accountID]
	if !ok {
		lock = &sync.Mutex{}
		m.refreshLocks[accountID] = lock
	}
	return lock
}

func (m *Manager) refreshAccount(account *model.CodexOAuthAccount) (string, error) {
	if account.RefreshToken == "" {
		return "", errors.New("codex oauth refresh token is missing")
	}
	return "", errors.New("codex oauth token refresh is not implemented yet")
}
