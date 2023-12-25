package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"log"

	"github.com/valli0x/booking-server/internal/cache"
	"github.com/valli0x/booking-server/internal/email"
	"github.com/valli0x/booking-server/internal/storage"
)

const (
	timeoutSeconds = 10
	idleTimeout    = 60
	maxHeaderBytes = 1024 * 1024
)

type server struct {
	// server fields
	addr string
	srv  *http.Server
	// services fields
	store  storage.Storage
	cache  cache.Cache
	mailer email.Mailer
}

type SrvConfig struct {
	Addr string
	// services fields
	Store  storage.Storage
	Cache  cache.Cache
	Mailer email.Mailer
}

func NewServer(cfg *SrvConfig) *server {
	httpServer := &http.Server{
		MaxHeaderBytes: maxHeaderBytes,
		IdleTimeout:    idleTimeout * time.Second,
		ReadTimeout:    timeoutSeconds * time.Second,
		WriteTimeout:   timeoutSeconds * time.Second,
	}

	s := &server{
		srv:    httpServer,
		addr:   cfg.Addr,
		store:  cfg.Store,
		cache:  cfg.Cache,
		mailer: cfg.Mailer,
	}

	s.srv.Handler = s.routers()

	return s
}

func (s *server) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("can't listen on %s. server quitting: %v", s.addr, err)
		return
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		<-ctx.Done()

		if err := s.srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}(wg)

	log.Printf("server listening on %s", s.addr)
	if err := s.srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("unexpected (http.Server).Serve error: %v", err)
	}

	wg.Wait()
	log.Printf("server off")
}
