package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"log"

	"github.com/go-chi/chi"
	cacheinmem "github.com/valli0x/booking-server/cache/inmem"
	"github.com/valli0x/booking-server/email"
	storeinmem "github.com/valli0x/booking-server/storage/inmem"
)

const (
	TimeoutSeconds = 10
	IdleTimeout    = 60
	MaxHeaderBytes = 1024 * 1024
)

type Server struct {
	// server fields
	addr string
	srv  *http.Server
	// services fields
	store  *storeinmem.InMemoryStorage
	cache  *cacheinmem.SimpleCache
	mailer *email.DummyMailer
}

type SrvConfig struct {
	Addr string
	// services fields
	store  *storeinmem.InMemoryStorage
	cache  *cacheinmem.SimpleCache
	mailer *email.DummyMailer
}

func NewServer(cfg *SrvConfig) (*Server, error) {
	httpServer := &http.Server{
		MaxHeaderBytes: MaxHeaderBytes,
		IdleTimeout:    IdleTimeout * time.Second,
		ReadTimeout:    TimeoutSeconds * time.Second,
		WriteTimeout:   TimeoutSeconds * time.Second,
	}

	s := &Server{
		srv:    httpServer,
		addr:   cfg.Addr,
		store:  cfg.store,
		cache:  cfg.cache,
		mailer: cfg.mailer,
	}

	s.srv.Handler = s.routers()

	return s, nil
}

func (s *Server) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("can't listen on %s. admin server quitting: %v", s.addr, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		if err := s.srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	log.Printf("server listening on %s", s.addr)
	if err := s.srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("unexpected (http.Server).Serve error: %v", err)
	}

	wg.Wait()
	log.Printf("server off")
}

func (s *Server) routers() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/orders", s.createOrder())
		r.Get("/orders/{userID}", s.getorders())
	})
	return r
}
