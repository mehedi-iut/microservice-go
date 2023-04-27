package main

import (
	"context"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"net/http"
	"os"
	"os/signal"
	"product-images/files"
	"product-images/handlers"
	"time"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9091", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
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

	// create CORS handler
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create the handlers
	fh := handlers.NewFiles(stor, l)
	mw := handlers.GzipHandler{}

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// filename regex: {filename:[a-zA-Z]+\\.[a-z]{3}}
	// problem with FileServer is that it is dumb
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadREST)
	ph.HandleFunc("/", fh.UploadMultipart)

	// get files
	gh := sm.Methods(http.MethodGet).Subrouter()
	// we are using "Handle" instead of HandleFunc as we are using default function http.FileServer
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))),
	)
	gh.Use(mw.GzipMiddleware)

	// create new server
	s := http.Server{
		Addr:         *bindAddress,
		Handler:      ch(sm),
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

	// Block until a signal is received
	sig := <-c
	l.Info("shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
