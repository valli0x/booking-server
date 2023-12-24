package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"log"

	cacheinmem "github.com/valli0x/booking-server/internal/cache/inmem"
	"github.com/valli0x/booking-server/internal/email"
	storeinmem "github.com/valli0x/booking-server/internal/storage/inmem"
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
	store  *storeinmem.InMemoryStorage
	cache  *cacheinmem.SimpleCache
	mailer *email.DummyMailer
}

type SrvConfig struct {
	Addr string
	// services fields
	Store  *storeinmem.InMemoryStorage
	Cache  *cacheinmem.SimpleCache
	Mailer *email.DummyMailer
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
