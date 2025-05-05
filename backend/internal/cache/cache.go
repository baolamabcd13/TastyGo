package cache

import (
    "sync"
    "time"
)

// Item đại diện cho một mục trong cache
type Item struct {
    Value      interface{}
    Expiration int64
}

// Cache là một bộ nhớ đệm đơn giản
type Cache struct {
    items map[string]Item
    mu    sync.RWMutex
}

// NewCache tạo một cache mới
func NewCache() *Cache {
    cache := &Cache{
        items: make(map[string]Item),
    }
    
    // Khởi động goroutine dọn dẹp các mục hết hạn
    go cache.janitor()
    
    return cache
}

// Set thêm một mục vào cache với thời gian hết hạn
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    var expiration int64
    if duration > 0 {
        expiration = time.Now().Add(duration).UnixNano()
    }
    
    c.items[key] = Item{
        Value:      value,
        Expiration: expiration,
    }
}

// Get lấy một mục từ cache
func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, found := c.items[key]
    if !found {
        return nil, false
    }
    
    // Kiểm tra hết hạn
    if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
        return nil, false
    }
    
    return item.Value, true
}

// Delete xóa một mục khỏi cache
func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    delete(c.items, key)
}

// Clear xóa tất cả các mục trong cache
func (c *Cache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.items = make(map[string]Item)
}

// janitor dọn dẹp các mục hết hạn
func (c *Cache) janitor() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.mu.Lock()
        now := time.Now().UnixNano()
        for key, item := range c.items {
            if item.Expiration > 0 && now > item.Expiration {
                delete(c.items, key)
            }
        }
        c.mu.Unlock()
    }
}

// Biến cache toàn cục
var DefaultCache = NewCache()

// Các hàm cache toàn cục
func Set(key string, value interface{}, duration time.Duration) {
    DefaultCache.Set(key, value, duration)
}

func Get(key string) (interface{}, bool) {
    return DefaultCache.Get(key)
}

func Delete(key string) {
    DefaultCache.Delete(key)
}

func Clear() {
    DefaultCache.Clear()
}