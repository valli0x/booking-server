package httpserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/valli0x/booking-server/pkg/models"
)

func (s *server) createOrder() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var o models.Order
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Error decoding order:", err)
			return
		}

		// Проверяем поля Order
		if err := s.validateOrder(&o); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("Error validate order:", err)
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
			orders, _ := s.store.GetOrders(rid)
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

		if err := s.store.Book(o); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error booking:", err)
			return
		}

		// Добавляем заказ в кеш
		if err := s.cache.AddOrder(o); err != nil {
			log.Println("Error adding order to cache:", err)
		}

		w.WriteHeader(http.StatusCreated)
		log.Println("Order created:", o)

		// Отправляем подтверждение
		if err := s.mailer.SendConfirmation(o.UserEmail, o); err != nil {
			log.Println("Error sending confirmation:", err)
		}
	})
}

func (s *server) validateOrder(order *models.Order) error {
	if order.UserID == "" || order.UserEmail == "" ||
		order.From == "" || order.To == "" || len(order.RoomIDs) == 0 || order.RoomType == "" {
		return errors.New("not all fields are filled")
	}

	_, err := mail.ParseAddress(order.UserEmail)
	if err != nil {
		return errors.New("invalid email")
	}

	_, err = time.Parse("2006-01-02", order.From)
	if err != nil {
		return errors.New("invalid 'From' date")
	}

	_, err = time.Parse("2006-01-02", order.To)
	if err != nil {
		return errors.New("invalid 'To' date")
	}

	return nil
}
