package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"net/http"
	"os"
	"os/signal"
	"product-image/files"
	"product-image/handlers"
	"time"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9091", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-image",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// max filesize 5MB
	stor, err := files.NewLocal(*basePath, 1024*1000*5)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// create the handlers
	fh := handlers.NewFiles(stor, l)
	mw := handlers.GzipHandler{}

	sm := chi.NewRouter()

	// Configure CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // You can specify specific origins instead of "*"
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	// Apply CORS middleware to the router
	sm.Use(corsHandler.Handler)

	sm.Post("/images/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadREST)
	sm.Post("/", fh.UploadMultipart)

	// get files
	sm.Group(func(r chi.Router) {
		r.Use(mw.GzipMiddleware)
		r.Get("/images/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/{filename:[a-zA-Z]+\\.[a-z]{3}}", func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))).ServeHTTP(w, r)
		})
	})

	// create new server
	s := http.Server{
		Addr:         *bindAddress,
		Handler:      sm,
		ErrorLog:     sl,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		l.Info("Starting the server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start the server")
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block untill a signal is received
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
