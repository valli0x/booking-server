package inmem

import (
	"sync"

	"github.com/valli0x/booking-server/pkg/models"
)

type InMemoryStorage struct {
	mu     sync.Mutex
	orders map[string][]models.Order
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		orders: make(map[string][]models.Order),
	}
}

func (s *InMemoryStorage) Book(o models.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[o.UserID] = append(s.orders[o.UserID], o)
	return nil
}

func (s *InMemoryStorage) GetOrders(userID string) ([]models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.orders[userID], nil
}
