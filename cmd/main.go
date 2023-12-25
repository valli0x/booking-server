package main

import (
	"context"
	"time"

	cacheinmem "github.com/valli0x/booking-server/internal/cache/inmem"
	"github.com/valli0x/booking-server/internal/email"
	"github.com/valli0x/booking-server/pkg/httpserver"
	storeinmem "github.com/valli0x/booking-server/internal/storage/inmem"
)

func main() {
	// ранее в проектах использовался пакет с di
	// но в данном случае мне нельзя использовать сторонние библиотеки
	store := storeinmem.NewInMemoryStorage()
	cache := cacheinmem.NewSimpleCache()
	mailer := email.NewDummyMailer()

	server := httpserver.NewServer(&httpserver.SrvConfig{
		Addr:   "localhost:8000",
		Store:  store,
		Cache:  cache,
		Mailer: mailer,
	})

	serverWorkTime := 10 * time.Minute // :)
	ctx, ctxcancel := context.WithTimeout(context.Background(), serverWorkTime)
	defer ctxcancel()

	server.Run(ctx)
}
