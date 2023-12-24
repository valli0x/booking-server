package inmem

import (
	"log"
	"sync"
	"time"

	"github.com/valli0x/booking-server/pkg/models"
)

type SimpleCache struct {
	mu     sync.Mutex
	orders map[string][]models.Order
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		orders: make(map[string][]models.Order),
	}
}

func (c *SimpleCache) AddOrder(o models.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, rid := range o.RoomIDs {
		c.orders[rid] = append(c.orders[rid], o)
	}
	return nil
}

func (c *SimpleCache) GetOrders(roomID string) ([]models.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.orders[roomID], nil
}

// Кеш, который учитывает время добавления ордера на бронь
// ? зачем я вообще это начал...
type Cache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value      models.Order
	expiration time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]cacheItem),
		mu:   sync.RWMutex{},
	}
}

func (c *Cache) Set(key string, value models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expirationTime, err := time.Parse("2006-01-02", value.To)
	if err != nil {
		log.Println("Failed to parse expiration time:", err)
		return
	}

	c.data[key] = cacheItem{
		value:      value,
		expiration: expirationTime,
	}
}

func (c *Cache) Get(key string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok {
		return models.Order{}, false
	}

	// Check if the item has expired
	if time.Now().After(item.expiration) {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.data, key)
		return models.Order{}, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *Cache) StartExpirationChecker() {
	go func() {
		for {
			time.Sleep(1 * time.Minute) // Check expiration every minute

			c.mu.Lock()
			for key, item := range c.data {
				if time.Now().After(item.expiration) {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
