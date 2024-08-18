package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	util "reservation_system/utilities"
)

// HandleDefault : Default API handler
/*
	Author		: Ravneet Singh
	Function 	: HandleDefault - Serves the default message for the default handler
	Input(s) 	: http.ResponseWriter, *http.Request (http.ResponseWriter to serve the response and *http.Request to process the request)
	Output(s) 	: None
	Example(s) 	: HandleDefault(response http.ResponseWriter, request *http.Request)
*/
func HandleDefault(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")
	fmt.Fprint(response, `{"message":"The server is up. Use other endpoints for reservation"}`)
}

// HandleAddReservation : Add reservation API Handler
/*
	Author		: Ravneet Singh
	Function 	: HandleAddReservation - Receives a POST request to add a reservation and adds it to a queue of reservations if no slots
				  until confirmed or cancelled
	Input(s) 	: http.ResponseWriter, *http.Request (http.ResponseWriter to serve the response and *http.Request to process the request)
	Output(s) 	: None
	Example(s) 	: HandleAddReservation(response http.ResponseWriter, request *http.Request)
*/
func HandleAddReservation(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != "POST" {
		http.Error(response, `{"error":"Wrong HTTP method specified"}`, http.StatusMethodNotAllowed)
		// fmt.Fprintf(response, util.ErrorMessageWrongMethod)
		return
	}

	reservationRequest := util.ReservationDetails{}
	err := json.NewDecoder(request.Body).Decode(&reservationRequest)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, util.ErrorMessageBadrequest)
		return
	}

	// Check for date time validity for the request
	isValid, err := IsRequestValid(reservationRequest)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, util.ErrorMessageBadrequest)
		return
	}

	if isValid {
		// Add the reservation to the queue or the datastore
		msg, err := AddReservation(reservationRequest)
		if err != nil {
			log.Println("Error in adding reservation: ", err)
			response.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(response, util.ErrorMessageServerError)
			return
		}
		log.Println("Request ", request.URL.Path, " processed successfully")
		response.WriteHeader(http.StatusOK)
		response.Header().Set("Content-Type", "application/json")
		fmt.Fprint(response, msg)
		return
	}
	log.Println("Not added, something went wrong")
	response.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(response, util.ErrorMessageServerError)
}

// HandleViewReservations : View reservation by date API Handler
/*
	Author		: Ravneet Singh
	Function 	: HandleViewReservations - Receives a GET request to view all reservations for a day time wise
	Input(s) 	: http.ResponseWriter, *http.Request (http.ResponseWriter to serve the response and *http.Request to process the request)
	Output(s) 	: None
	Example(s) 	: HandleViewReservations(response http.ResponseWriter, request *http.Request)
*/
func HandleViewReservations(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != "GET" {
		http.Error(response, `{"error":"Wrong HTTP method specified"}`, http.StatusMethodNotAllowed)
		fmt.Fprintf(response, util.ErrorMessageWrongMethod)
		return
	}

	param := request.URL.Query().Get("date")
	if len(param) != 10 {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, util.ErrorMessageBadrequest)
		return
	}

	result, err := ViewReservationsByDate("default", param)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, util.ErrorMessageServerError)
		return
	}

	jsonResponse, err := json.Marshal(result)
	if err != nil {
		log.Println("Something went wrong. ", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, util.ErrorMessageServerError)
	} else {
		log.Println("Request ", request.URL.Path, "processed successfully!")
		response.WriteHeader(http.StatusOK)
		fmt.Fprintf(response, "%s", string(jsonResponse[:]))
	}
}

// HandleCancelReservation : Cancel reservation API Handler
/*
	Author		: Ravneet Singh
	Function 	: HandleCancelReservation - Receives a POST request to cancel a reservation and remove it from the datastore
	Input(s) 	: http.ResponseWriter, *http.Request (http.ResponseWriter to serve the response and *http.Request to process the request)
	Output(s) 	: None
	Example(s) 	: HandleCancelReservation(response http.ResponseWriter, request *http.Request)
*/
func HandleCancelReservation(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != "POST" {
		http.Error(response, `{"error":"Wrong HTTP method specified"}`, http.StatusMethodNotAllowed)
		fmt.Fprintf(response, util.ErrorMessageWrongMethod)
		return
	}

	cancellationrequest := util.ReservationDetails{}
	err := json.NewDecoder(request.Body).Decode(&cancellationrequest)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, util.ErrorMessageBadrequest)
		return
	}

	result, err := CancelReservation(cancellationrequest)
	if err != nil {
		log.Println("Error in cancelation of request: ", err)
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, util.ErrorMessageServerError)
		return
	}

	log.Println("Request ", request.URL.Path, "processed successfully!")
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")
	fmt.Fprint(response, result)
}

// HandleConfirmReservation : Confirm reservation API Handler
/*
	Author		: Ravneet Singh
	Function 	: HandleConfirmReservation - Receives a POST request to confirm a reservation and adds it to a the datastore
	Input(s) 	: http.ResponseWriter, *http.Request (http.ResponseWriter to serve the response and *http.Request to process the request)
	Output(s) 	: None
	Example(s) 	: HandleConfirmReservation(response http.ResponseWriter, request *http.Request)
*/
func HandleConfirmReservation(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != "POST" {
		http.Error(response, `{"error":"Wrong HTTP method specified"}`, http.StatusMethodNotAllowed)
		fmt.Fprintf(response, util.ErrorMessageWrongMethod)
		return
	}

	confirmation := util.ConfirmReservation{}
	err := json.NewDecoder(request.Body).Decode(&confirmation)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, util.ErrorMessageBadrequest)
		return
	}

	msg, err := ConfirmReservation(confirmation)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, util.ErrorMessageServerError)
		return
	}

	log.Println("Request ", request.URL.Path, "processed successfully!")
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")
	fmt.Fprint(response, msg)
}

// HandleViewWaitingList : View Waiting List of reservation API Handler
/*
	Author		: Ravneet Singh
	Function 	: HandleViewWaitingList - Receives a GET request to view all reservations for a day time wise
	Input(s) 	: http.ResponseWriter, *http.Request (http.ResponseWriter to serve the response and *http.Request to process the request)
	Output(s) 	: None
	Example(s) 	: HandleViewWaitingList(response http.ResponseWriter, request *http.Request)
*/
func HandleViewWaitingList(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != "GET" {
		http.Error(response, `{"error":"Wrong HTTP method specified"}`, http.StatusMethodNotAllowed)
		fmt.Fprintf(response, util.ErrorMessageWrongMethod)
		return
	}

	queue.Mutex.Lock()
	result := queue.Queue
	queue.Mutex.Unlock()
	jsonResponse, err := json.Marshal(result)
	if err != nil {
		log.Println("Something went wrong. ", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, util.ErrorMessageServerError)
	}

	log.Println("Request ", request.URL.Path, "processed successfully!")
	response.WriteHeader(http.StatusOK)
	fmt.Fprintf(response, "%s", string(jsonResponse[:]))
}

// HandleShowAvailability : Show available slots by location and day
/*
	Author		: Ravneet Singh
	Function 	: HandleShowAvailability - Receives a GET request to view all reservations for a day time wise
	Input(s) 	: http.ResponseWriter, *http.Request (http.ResponseWriter to serve the response and *http.Request to process the request)
	Output(s) 	: None
	Example(s) 	: HandleShowAvailability(response http.ResponseWriter, request *http.Request)
*/
func HandleShowAvailability(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != "GET" {
		http.Error(response, `{"error":"Wrong HTTP method specified"}`, http.StatusMethodNotAllowed)
		fmt.Fprintf(response, util.ErrorMessageWrongMethod)
		return
	}

	location := request.URL.Query().Get("location")
	if location == "" {
		location = "default"
	}
	date := request.URL.Query().Get("date")
	if len(date) != 10 {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, util.ErrorMessageBadrequest)
		return
	}

	result := GetAvailabilities(location, date)
	jsonResponse, err := json.Marshal(result)
	if err != nil {
		log.Println("Something went wrong. ", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(response, util.ErrorMessageServerError)
	}
	log.Println("Request ", request.URL.Path, "processed successfully!")
	response.WriteHeader(http.StatusOK)
	fmt.Fprintf(response, "%s", string(jsonResponse[:]))
}
