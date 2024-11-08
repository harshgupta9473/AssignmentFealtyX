package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/harshgupta9473/GoAssignmentFealtyX/handlers"
	"github.com/harshgupta9473/GoAssignmentFealtyX/routes"
)

func main() {
	studentHandler := handlers.NewStudentHandler()
	router := mux.NewRouter()
	routes.RegisterRoutes(router, studentHandler)

	s := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		// IdleTimeout:  120 * time.Second,
		// ReadTimeout:  1 * time.Second,
		// WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan:=make(chan os.Signal)
	signal.Notify(sigChan,os.Interrupt)
	signal.Notify(sigChan,os.Kill)

	sig:=<-sigChan
	log.Println("recieved terminate signal, graceful shutdown",sig)

	tc,_:=context.WithTimeout(context.Background(),30*time.Second)
	s.Shutdown(tc)
}
