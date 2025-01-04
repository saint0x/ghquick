package cache

import (
	"sync"
	"time"
)

type RepoInfo struct {
	Name      string
	Path      string
	Remote    string
	Branch    string
	UpdatedAt time.Time
}

type RepoCache struct {
	mu    sync.RWMutex
	cache map[string]*RepoInfo
	ttl   time.Duration
}

func NewRepoCache() *RepoCache {
	return &RepoCache{
		cache: make(map[string]*RepoInfo),
		ttl:   5 * time.Minute, // Cache TTL of 5 minutes
	}
}

func (c *RepoCache) Get(path string) (*RepoInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	info, exists := c.cache[path]
	if !exists {
		return nil, false
	}

	if time.Since(info.UpdatedAt) > c.ttl {
		delete(c.cache, path)
		return nil, false
	}

	return info, true
}

func (c *RepoCache) Set(path string, info *RepoInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	info.UpdatedAt = time.Now()
	c.cache[path] = info
}
