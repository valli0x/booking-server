package storage

import "github.com/valli0x/booking-server/models"

type Storage interface {
	Book(o models.Order) error
	GetOrders(userID string) ([]models.Order, error)
}
