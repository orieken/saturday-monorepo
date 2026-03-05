package templates

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// CacheEntry represents a cached template result
type CacheEntry struct {
	Result    string
	Timestamp time.Time
}

// Cache provides caching for processed templates
type Cache struct {
	entries map[string]CacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
}

// NewCache creates a new template cache
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
	}
}

// Get retrieves a cached template result
func (c *Cache) Get(templateName string, data interface{}) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.generateKey(templateName, data)
	entry, exists := c.entries[key]

	if !exists {
		return "", false
	}

	// Check if entry has expired
	if c.ttl > 0 && time.Since(entry.Timestamp) > c.ttl {
		return "", false
	}

	return entry.Result, true
}

// Set stores a template result in the cache
func (c *Cache) Set(templateName string, data interface{}, result string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.generateKey(templateName, data)
	c.entries[key] = CacheEntry{
		Result:    result,
		Timestamp: time.Now(),
	}
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]CacheEntry)
}

// Size returns the number of cached entries
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.entries)
}

// Cleanup removes expired entries from the cache
func (c *Cache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ttl == 0 {
		return 0
	}

	removed := 0
	now := time.Now()

	for key, entry := range c.entries {
		if now.Sub(entry.Timestamp) > c.ttl {
			delete(c.entries, key)
			removed++
		}
	}

	return removed
}

// generateKey creates a cache key from template name and data
func (c *Cache) generateKey(templateName string, data interface{}) string {
	// Serialize data to JSON for consistent hashing
	dataBytes, err := json.Marshal(data)
	if err != nil {
		// Fallback to template name only if serialization fails
		return templateName
	}

	// Create hash of template name + data
	hash := sha256.Sum256(append([]byte(templateName), dataBytes...))
	return fmt.Sprintf("%x", hash)
}
