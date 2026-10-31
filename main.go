package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/coalaura/plain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var log = plain.New(plain.WithDate(plain.RFC3339Local))

func main() {
	var err error

	log.Println("Loading config...")

	config, err := LoadConfig()
	log.MustFail(err)

	log.Println("Preparing cache...")

	cache := NewCache()

	log.Println("Preparing router...")

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(plain.RealIPMiddleware())
	r.Use(log.Middleware())

	handler := HandleRequest(config, cache)

	r.Get("/{name}", handler)
	r.Get("/{name}.{format}", handler)

	addr := config.Addr()

	server := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Listening on http://localhost%s\n", addr)

		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warnln(err)
		}
	}()

	log.WaitForInterrupt()

	log.Warnln("Shutting down...")

	server.Close()
}
