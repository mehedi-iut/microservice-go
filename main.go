package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"practice/data"
	"practice/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	connectionString := "server=bs054.database.windows.net;user id=mehedi;password=database@123;port=1433;database=test"
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	p := &data.ProductModel{DB: db}
	pl := &data.ProductInfo{}

	ph := handlers.NewProducts(l, p, pl)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}

	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = s.Shutdown(tc)

	if err != nil {
		l.Fatalf("HTTP server shutdown, %v", err)
	}

}
