package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	server, err := NewConfigServer()
	if err != nil {
		log.Fatal(err)
		return
	}

	router.HandleFunc("/config/", count(server.createConfHandler)).Methods("POST")
	router.HandleFunc("/group/", count(server.createConfGroupHandler)).Methods("POST")
	router.HandleFunc("/configs/", count(server.getConfigsHandler)).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/", count(server.delConfigHandler)).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/", count(server.delConfigGroupsHandler)).Methods("DELETE")
	router.HandleFunc("/config/{id}/{version}/", count(server.viewConfigHandler)).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/", count(server.viewGroupHandler)).Methods("GET")
	// router.HandleFunc("/group/{id}/{version}/{label}/", count(server.viewGroupLabelHandler)).Methods("GET")
	// router.HandleFunc("/group/{id}/{version}/", count(server.updateConfigHandler)).Methods("PUT")
	router.Path("/metrics").Handler(metricsHandler())

	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
