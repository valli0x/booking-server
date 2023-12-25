package cache

import (
	"github.com/valli0x/booking-server/pkg/models"
)

// Можно было бы создать абстрактный кеш, который обычно и используется
// по типу lru кеша, где значение interface, то как то так
// микросервис работает в основном с ордерами(заказами)
type Cache interface {
	AddOrder(o models.Order) error
	GetOrders(roomID string) ([]models.Order, error)
}
