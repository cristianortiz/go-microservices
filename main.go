package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cristianortiz/go-microservices/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//read .env files in rootPath TODO: checks if .env file exists before
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file does not exist")
	}
	//Bind address for the server
	port := os.Getenv("SERVER_PORT")

	l := log.New(os.Stdout, "api ", log.LstdFlags)

	// create the handlers
	ph := handlers.NewProducts(l)

	// create a new server with gorilla mux
	sm := mux.NewRouter()
	//config a subrouter to filter routes to a particular httpp method and http handler
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	//specific handlerfunc GetProduct associated with route "/"
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	//use regexp to ensures id param contains only digits
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProducts)
	postRouter.Use(ph.MiddlewareProductValidation)

	// create a new server
	s := http.Server{
		Addr:         port,              // configure port
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server using goroutine
	go func() {
		l.Println("Starting server on port", port)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
	//explicit cancel() in context.WithTimeout and defer it later to avoid warning about cancel function not discarded
	defer cancel()

}
