package cache

import (
	"github.com/valli0x/booking-server/pkg/models"
)

type Cache interface {
	AddOrder(o models.Order) error
	GetOrders(roomID string) ([]models.Order, error)
}
