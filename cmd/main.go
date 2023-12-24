package main

import (
	"context"
	"time"

	cacheinmem "github.com/valli0x/booking-server/cache/inmem"
	"github.com/valli0x/booking-server/email"
	"github.com/valli0x/booking-server/httpserver"
	storeinmem "github.com/valli0x/booking-server/storage/inmem"
)

func main() {
	store := storeinmem.NewInMemoryStorage()
	cache := cacheinmem.NewSimpleCache()
	mailer := email.NewDummyMailer()

	server := httpserver.NewServer(&httpserver.SrvConfig{
		Addr: ":3000",
		Store: store,
		Cache: cache,
		Mailer: mailer,
	})

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Minute)
	go server.Run(ctx)
}
