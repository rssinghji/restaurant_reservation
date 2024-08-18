package main

import (
	"log"
	"net/http"
	"reservation_system/api"
	"runtime"
	"time"
)

func main() {
	log.Println("Initializing restaurant reservation system........")

	// Set the maximum number of CPUs that can be executing simultaneously
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a server with optimized settings
	server := &http.Server{
		Addr:         ":7070",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// All the handlers for the reservation system apis
	http.HandleFunc("/", api.HandleDefault)
	http.HandleFunc("/add", api.HandleAddReservation)
	http.HandleFunc("/confirm", api.HandleConfirmReservation)
	http.HandleFunc("/view", api.HandleViewReservations)
	http.HandleFunc("/cancel", api.HandleCancelReservation)
	http.HandleFunc("/waitinglist", api.HandleViewWaitingList)
	http.HandleFunc("/availability", api.HandleShowAvailability)

	// Print the startup message
	log.Println("Starting server on :7070")
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server failed to start: %v\n", err)
	}
}
