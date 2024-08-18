package utilities

import (
	"sync"
	"time"
)

// ReservationDetails: Basic details of reservation
type ReservationDetails struct {
	Name     string `json:"name"`     // name of the requestor
	Date     string `json:"date"`     // in yyyy-mm-dd format
	Time     string `json:"time"`     // in HH:MM AM format
	Email    string `json:"email"`    // Send notificaiton to confirm reservation
	Location string `json:"location"` // use "default" if not mentioned
}

type ConfirmReservation struct {
	Name     string `json:"name"`     // name of the requestor
	Date     string `json:"date"`     // in yyyy-mm-dd format
	Time     string `json:"time"`     // in HH:MM AM format
	Location string `json:"location"` // use "default" if not mentioned
}

type ReservationResponse struct {
	Name        string    `json:"name"`
	Date        string    `json:"date"`     // in yyyy-mm-dd format
	Time        string    `json:"time"`     // in HH:MM AM format
	Email       string    `json:"email"`    // send notificaiton to confirm reservation
	Status      string    `json:"status"`   // either CONFIRMED or WAITING, CANCELED ones would be deleted
	Location    string    `json:"location"` // use "default" if not mentioned
	RequestTime time.Time `json:"-"`        // for internal purposes only
}

// Reservation : Strutcture to store everything about reservation details
type Reservation struct {
	Mutex      sync.Mutex
	CountMap   map[string]int
	DetailsMap map[string]ReservationResponse
}

// ReservationQueue : Queue to store incoming requests and serve them in FIFO order
type ReservationQueue struct {
	Mutex   sync.Mutex
	Windows []string
	Queue   []ReservationDetails
}

// Constants for API handling
const (
	ErrorMessageWrongMethod      = `{"error":"Wrong HTTP method."}`
	ErrorMessageBadrequest       = `{"error":"Bad Request."}`
	ErrorMessageServerError      = `{"error":"Something went wrong."}`
	ErrorConfirmReservation      = `{"error":"Reservation doesn't exist. Probably timed out. Make a new one."}`
	SuccessAddReservationMessage = `{"message":"Added your reservation, please confirm within 2 minutes"}`
	SuccessAddWaitlistMessage    = `{"message":"Added you to waitlist, would notify if available"}`
	SuccessCancellationMessage   = `{"message":"Your reservation has been canceled upon your request"}`
	SuccessConfirmReservation    = `{"message":"Your reservation has been confirmed upon your request"}`
)
