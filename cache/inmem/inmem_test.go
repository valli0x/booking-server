package inmem

import (
	"testing"
	"time"

	"github.com/valli0x/booking-server/models"
)

func TestInmem(t *testing.T) {
	cache := NewCache()
	cache.StartExpirationChecker()

	order := models.Order{
		UserID:    "123",
		UserEmail: "john@example.com",
		RoomType:  models.Econom,
		RoomIDs:   []string{"1", "2"},
		From:      "2022-01-01", // TODO: adding hour
		To:        "2022-01-01", // TODO: adding hour
	}

	cache.Set("order123", order)

	cachedOrder, ok := cache.Get("order123")
	if !ok {
		t.Fatal("Order found in cache, Order:", cachedOrder)
	}

	time.Sleep(1 * time.Hour)

	_, ok = cache.Get("order123")
	if ok {
		t.Fatal("Order found in cache")
	}
}
