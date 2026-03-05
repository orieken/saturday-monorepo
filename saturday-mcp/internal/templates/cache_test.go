package templates

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		data := map[string]string{"name": "test"}
		cache.Set("template1", data, "result1")

		result, found := cache.Get("template1", data)
		if !found {
			t.Error("Expected to find cached result")
		}

		if result != "result1" {
			t.Errorf("Expected result1, got %s", result)
		}
	})

	t.Run("Get non-existent", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		_, found := cache.Get("nonexistent", nil)
		if found {
			t.Error("Expected not to find non-existent cache entry")
		}
	})

	t.Run("TTL expiration", func(t *testing.T) {
		cache := NewCache(100 * time.Millisecond)

		data := map[string]string{"name": "test"}
		cache.Set("template1", data, "result1")

		// Should be found immediately
		_, found := cache.Get("template1", data)
		if !found {
			t.Error("Expected to find cached result before expiration")
		}

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		_, found = cache.Get("template1", data)
		if found {
			t.Error("Expected cache entry to be expired")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		cache.Set("template1", nil, "result1")
		cache.Set("template2", nil, "result2")

		if cache.Size() != 2 {
			t.Errorf("Expected size 2, got %d", cache.Size())
		}

		cache.Clear()

		if cache.Size() != 0 {
			t.Errorf("Expected size 0 after clear, got %d", cache.Size())
		}
	})

	t.Run("Cleanup", func(t *testing.T) {
		cache := NewCache(100 * time.Millisecond)

		cache.Set("template1", nil, "result1")
		cache.Set("template2", nil, "result2")

		time.Sleep(150 * time.Millisecond)

		removed := cache.Cleanup()
		if removed != 2 {
			t.Errorf("Expected to remove 2 entries, removed %d", removed)
		}

		if cache.Size() != 0 {
			t.Errorf("Expected size 0 after cleanup, got %d", cache.Size())
		}
	})

	t.Run("Different data generates different keys", func(t *testing.T) {
		cache := NewCache(5 * time.Minute)

		data1 := map[string]string{"name": "test1"}
		data2 := map[string]string{"name": "test2"}

		cache.Set("template1", data1, "result1")
		cache.Set("template1", data2, "result2")

		result1, _ := cache.Get("template1", data1)
		result2, _ := cache.Get("template1", data2)

		if result1 != "result1" {
			t.Errorf("Expected result1, got %s", result1)
		}

		if result2 != "result2" {
			t.Errorf("Expected result2, got %s", result2)
		}
	})
}
