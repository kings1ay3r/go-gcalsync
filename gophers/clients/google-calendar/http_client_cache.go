package googlecalendar

import (
	"net/http"
	"sync"
	"time"
)

type Client struct {
	httpClient *http.Client
	createdAt  time.Time
}

type ClientCache struct {
	mu       sync.RWMutex
	clients  map[string]*Client
	expiry   time.Time
	lifetime time.Duration
}

func NewClientCache(lifetime time.Duration) *ClientCache {
	return &ClientCache{
		clients:  make(map[string]*Client),
		expiry:   time.Now().Add(lifetime),
		lifetime: lifetime,
	}
}

func (c *ClientCache) Push(userID string, client *http.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.expiry.Before(time.Now()) {
		c.reset()
	}
	c.clients[userID] = &Client{
		httpClient: client,
		createdAt:  time.Now(),
	}
}

func (c *ClientCache) Get(userID string) (*http.Client, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.expiry.Before(time.Now()) {
		c.reset()
		return nil, false
	}

	client, exists := c.clients[userID]
	if !exists {
		return nil, false
	}
	return client.httpClient, true
}

func (c *ClientCache) reset() {
	c.clients = make(map[string]*Client)
	c.expiry = time.Now().Add(c.lifetime)
}
