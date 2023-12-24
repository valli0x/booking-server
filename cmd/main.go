package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	cacheinmem "github.com/valli0x/booking-server/cache/inmem"
	"github.com/valli0x/booking-server/email"
	"github.com/valli0x/booking-server/models"
	storeinmem "github.com/valli0x/booking-server/storage/inmem"
)

var (
	store  = storeinmem.NewInMemoryStorage()
	cache  = cacheinmem.NewSimpleCache()
	mailer = email.NewDummyMailer()
)

func main() {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
			var o models.Order
			if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("Error decoding order:", err)
				return
			}

			// Проверяем тип комнаты
			if o.RoomType != models.Econom && o.RoomType != models.Standart && o.RoomType != models.Lux {
				http.Error(w, "Invalid room type", http.StatusBadRequest)
				log.Println("Invalid room type:", o.RoomType)
				return
			}

			// Проверка на пересечение времени бронирования для каждой комнаты
			for _, rid := range o.RoomIDs {
				orders, _ := cache.GetOrders(rid)
				for _, order := range orders {
					fromTime, _ := time.Parse(time.RFC3339, order.From)
					toTime, _ := time.Parse(time.RFC3339, order.To)
					newFromTime, _ := time.Parse(time.RFC3339, o.From)
					newToTime, _ := time.Parse(time.RFC3339, o.To)

					if fromTime.Before(newToTime) && newFromTime.Before(toTime) {
						http.Error(w, "Booking times conflict", http.StatusBadRequest)
						log.Println("Booking times conflict for order:", o)
						return
					}
				}
			}

			if err := store.Book(o); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error booking:", err)
				return
			}

			// Добавляем заказ в кеш
			if err := cache.AddOrder(o); err != nil {
				log.Println("Error adding order to cache:", err)
			}

			w.WriteHeader(http.StatusCreated)
			log.Println("Order created:", o)

			// Отправляем подтверждение
			if err := mailer.SendConfirmation(o.UserEmail, o); err != nil {
				log.Println("Error sending confirmation:", err)
			}
		})

		r.Get("/orders/{userID}", func(w http.ResponseWriter, r *http.Request) {
			userID := chi.URLParam(r, "userID")
			orders, err := store.GetOrders(userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error getting orders:", err)
				return
			}
			if err := json.NewEncoder(w).Encode(orders); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error encoding orders:", err)
				return
			}
			log.Println("Orders retrieved for user:", userID)
		})
	})

	log.Println("Starting server on :3000")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
