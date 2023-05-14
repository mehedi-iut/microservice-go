package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"practice/data"
	"practice/handlers"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

func main() {
	// l := log.New(os.Stdout, "product-api", log.LstdFlags)
	l := hclog.Default()

	err := godotenv.Load(".env")
	if err != nil {
		l.Error("Error getting env variable", "error", err)
	}

	var server = os.Getenv("server")
	var strPort = os.Getenv("port")
	var user = os.Getenv("user")
	var password = os.Getenv("password")
	var database = os.Getenv("database")

	port, err := strconv.Atoi(strPort)
	if err != nil {
		l.Error("Can't convert the port into integer", "Error", err)
	}
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		l.Error("Error loading Env variable", "error", err)
	}
	defer db.Close()

	p := &data.ProductModel{DB: db}
	pl := &data.ProductInfo{}

	ph := handlers.NewProducts(l, p, pl)

	//sm := http.NewServeMux()
	//sm.Handle("/", ph)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProducts)
	postRouter.Use(ph.MiddlewareValidateProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{name}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Info("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error Starting Server", "error", err)
		}

	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = s.Shutdown(tc)

	if err != nil {
		l.Error("HTTP server shutdown", "error", err)
	}

}
