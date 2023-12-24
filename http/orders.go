package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) getorders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		orders, err := s.store.GetOrders(userID)
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
	}
}
