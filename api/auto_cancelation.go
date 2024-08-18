package api

import (
	"log"
	"time"
)

// waitForConfirmation() : A forver running go routine which will check for new requests and keep track of confirmations
func waitForConfirmation(windowKey, detailsKey string, seconds int) {
	log.Println("Starting the confirmation goroutine")
	timeNow := time.Now()
	timeLater := timeNow.Add(time.Duration(time.Duration(seconds).Seconds())) // Change it to different time, if needed
	for timeNow.Before(timeLater) {
		reservation.Mutex.Lock()
		if reservation.DetailsMap[detailsKey].Status == "CONFIRMED" {
			reservation.Mutex.Unlock()
			return
		}
		reservation.Mutex.Unlock()
		time.Sleep(time.Minute)
		timeNow = timeNow.Add(time.Minute)
	}

	// Couldn't confirm in due time, so delete the entry
	log.Println("Auto deleting the request after timeout for: ", detailsKey)
	reservation.Mutex.Lock()
	delete(reservation.DetailsMap, detailsKey)
	count := reservation.CountMap[windowKey]
	reservation.CountMap[windowKey] = count - 1
	reservation.Mutex.Unlock()
}
