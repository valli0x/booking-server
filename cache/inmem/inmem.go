package inmem

import (
	"sync"

	"github.com/valli0x/booking-server/models"
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
