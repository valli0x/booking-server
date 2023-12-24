package storage

import "github.com/valli0x/booking-server/pkg/models"

type Storage interface {
	Book(o models.Order) error
	GetOrders(userID string) ([]models.Order, error)
}
