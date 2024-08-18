package api

import (
	"fmt"
	"log"
	util "reservation_system/utilities"
	"sort"
	"strconv"
	"strings"
	"time"
)

var queue util.ReservationQueue
var reservation util.Reservation
var windowOrder = []string{"10AM-11AM", "11AM-12PM", "12PM-1PM", "1PM-2PM", "2PM-3PM", "3PM-4PM", "4PM-5PM", "5PM-6PM", "6PM-7PM",
	"7PM-8PM", "8PM-9PM", "9PM-10PM"}
var peakTimes = []string{"6PM-7PM", "7PM-8PM"}

// init() : initialize the queue and reservation system
func init() {
	windowsList := make([]string, 0)
	reservationList := make([]util.ReservationDetails, 0)
	queue.Queue = reservationList
	queue.Windows = windowsList

	reservation.CountMap = make(map[string]int)
	reservation.DetailsMap = make(map[string]util.ReservationResponse)
	go waitListCleanUp(1800)
}

// getWindow() : A function to convert request to a time window
func getWindow(time string) (string, error) {
	log.Println("Running getWindow()")
	timeOfDay := "AM" // default it to morning
	timeSplits := strings.Split(time, ":")
	minutes := strings.Split(timeSplits[1], " ")
	timeOfDay = minutes[1]
	timeHour, err := strconv.Atoi(timeSplits[0])
	if err != nil {
		log.Println("Error in hour extraction for reservation")
		return "Error in hour extraction for reservation", fmt.Errorf("error in hour extraction for reservation")
	}

	switch {
	case timeHour >= 10 && timeHour < 11 && timeOfDay == "AM":
		return "10AM-11AM", nil
	case timeHour >= 11 && timeHour < 12 && timeOfDay == "AM":
		return "11AM-12PM", nil
	case timeHour >= 12 && timeOfDay == "PM":
		return "12PM-1PM", nil
	case timeHour >= 1 && timeHour < 2 && timeOfDay == "PM":
		return "1PM-2PM", nil
	case timeHour >= 2 && timeHour < 3 && timeOfDay == "PM":
		return "2PM-3PM", nil
	case timeHour >= 3 && timeHour < 4 && timeOfDay == "PM":
		return "3PM-4PM", nil
	case timeHour >= 4 && timeHour < 5 && timeOfDay == "PM":
		return "4PM-5PM", nil
	case timeHour >= 5 && timeHour < 6 && timeOfDay == "PM":
		return "5PM-6PM", nil
	case timeHour >= 6 && timeHour < 7 && timeOfDay == "PM":
		return "6PM-7PM", nil
	case timeHour >= 7 && timeHour < 8 && timeOfDay == "PM":
		return "7PM-8PM", nil
	case timeHour >= 8 && timeHour < 9 && timeOfDay == "PM":
		return "8PM-9PM", nil
	case timeHour >= 9 && timeHour < 10 && timeOfDay == "PM":
		return "9PM-10PM", nil
	default:
		// Should never come here
		return "Requested time is outside working hours, open from 10 AM to 10 PM", fmt.Errorf("requested time is outside working hours, open from 10 AM to 10 PM")
	}
}

// getMaxOccupancy() : A function to tell how many open slots based on peak times
func getMaxOccupancy(window string) int {
	maxOccupancy := 4
	for _, value := range peakTimes {
		if window == value {
			maxOccupancy = 3
			break
		}
	}
	return maxOccupancy
}

// IsRequestValid() : To check if the reservation request is valid against time constraints
func IsRequestValid(resRequest util.ReservationDetails) (bool, error) {
	log.Println("Running IsRequestValid()")
	currentDateTime := time.Now().Local()
	layout := "2006-01-02 15:04 PM"
	current := currentDateTime.Format(layout)
	reqDateTime := resRequest.Date + " " + resRequest.Time
	requestedTime, err := time.Parse(layout, reqDateTime)
	if err != nil {
		return false, fmt.Errorf("error in parsing date-time from requested value")
	}
	originalTime, _ := time.Parse(layout, current)

	// If requestedTime has already passed
	if requestedTime.Before(originalTime) {
		return false, fmt.Errorf("requested time has already passed, try new time and at least 5 minutes ahead")
	}

	if requestedTime.Hour() < 10 || requestedTime.Hour() >= 22 {
		return false, fmt.Errorf("requested time is outside working hours, open from 10 AM to 10 PM")
	}

	if requestedTime.After(originalTime.Add(5 * time.Minute)) {
		return true, nil
	}
	return false, nil // internal server error, something happened and code reached here
}

// AddReservation() : Adds a reservation to the datatore
func AddReservation(resReq util.ReservationDetails) (string, error) {
	log.Println("Running addRequest()")
	responseObject := util.ReservationResponse{Name: resReq.Name, Date: resReq.Date, Time: resReq.Time, Email: resReq.Email}
	location := "default"
	if resReq.Location != "" {
		location = resReq.Location
	}
	responseObject.Location = location
	responseObject.Status = "WAITING"
	window, err := getWindow(resReq.Time)
	if err != nil {
		return err.Error(), err
	}

	// Fetch max occupancy for the reservation window
	maxOccupancy := getMaxOccupancy(window)

	// Prepare to store info
	windowKey := location + "-" + resReq.Date + "-" + window
	detailsKey := windowKey + "-" + resReq.Name
	msg := util.SuccessAddReservationMessage
	responseObject.RequestTime = time.Now()
	if count, ok := reservation.CountMap[windowKey]; ok {
		if count == maxOccupancy {
			queue.Mutex.Lock()
			queue.Windows = append(queue.Windows, window)
			queue.Queue = append(queue.Queue, resReq)
			queue.Mutex.Unlock()
			msg = util.SuccessAddWaitlistMessage
			return msg, nil
		} else {
			reservation.Mutex.Lock()
			reservation.DetailsMap[detailsKey] = responseObject
			reservation.CountMap[windowKey] = count + 1
			reservation.Mutex.Unlock()
		}
	} else {
		reservation.Mutex.Lock()
		reservation.DetailsMap[detailsKey] = responseObject
		reservation.CountMap[windowKey] = 1
		reservation.Mutex.Unlock()
	}

	// Launch a go routine for auto cancellation in 2 minutes or 120 seconds
	go waitForConfirmation(windowKey, detailsKey, 120)

	return msg, nil
}

// ViewReservationsByDate() : Returns all the reservations by date for a location
func ViewReservationsByDate(location, date string) (map[string][]util.ReservationResponse, error) {
	result := make(map[string][]util.ReservationResponse)
	reservation.Mutex.Lock()
	defer reservation.Mutex.Unlock()
	for _, window := range windowOrder {
		reservedSlots := make([]util.ReservationResponse, 0)
		for key, value := range reservation.DetailsMap {
			compositeKeys := strings.Split(key, "-")
			dateMatch := compositeKeys[1] + "-" + compositeKeys[2] + "-" + compositeKeys[3]
			windowMatch := compositeKeys[4] + "-" + compositeKeys[5]
			if (location == compositeKeys[0]) && (date == dateMatch) &&
				(window == windowMatch) {
				reservedSlots = append(reservedSlots, value)
			}
		}
		sort.Slice(reservedSlots, func(i, j int) bool {
			return reservedSlots[i].RequestTime.Before(reservedSlots[j].RequestTime)
		})
		result[window] = reservedSlots
	}

	return result, nil
}

// CancelReservation() : Cancels a reservation and removes it from a datastore
func CancelReservation(resReq util.ReservationDetails) (string, error) {
	reservation.Mutex.Lock()
	location := "default"
	if resReq.Location != "" {
		location = resReq.Location
	}
	window, err := getWindow(resReq.Time)
	if err != nil {
		reservation.Mutex.Unlock()
		return err.Error(), err
	}
	windowKey := location + "-" + resReq.Date + "-" + window
	detailsKey := windowKey + "-" + resReq.Name

	// Remove the reservation and reduce the count
	delete(reservation.DetailsMap, detailsKey)
	count := reservation.CountMap[windowKey]
	reservation.CountMap[windowKey] = count - 1
	reservation.Mutex.Unlock()
	go manageWaitlist(window)
	return util.SuccessCancellationMessage, nil
}

// ConfirmReservation() : Confirms a reservation and change the status to CONFIRMED
func ConfirmReservation(confirmRequest util.ConfirmReservation) (string, error) {
	reservation.Mutex.Lock()
	defer reservation.Mutex.Unlock()
	location := "default"
	if confirmRequest.Location != "" {
		location = confirmRequest.Location
	}
	window, err := getWindow(confirmRequest.Time)
	if err != nil {
		return util.ErrorConfirmReservation, err
	}
	windowKey := location + "-" + confirmRequest.Date + "-" + window
	detailsKey := windowKey + "-" + confirmRequest.Name

	if value, ok := reservation.DetailsMap[detailsKey]; ok {
		value.Status = "CONFIRMED"
		reservation.DetailsMap[detailsKey] = value
	} else {
		return util.ErrorConfirmReservation, nil
	}

	return util.SuccessConfirmReservation, nil
}

// GetAvailabilities() : Function to return available spots for reservation
func GetAvailabilities(location, date string) map[string]int {
	result := make(map[string]int)

	reservation.Mutex.Lock()
	defer reservation.Mutex.Unlock()
	for _, value := range windowOrder {
		maxOccupancy := getMaxOccupancy(value)
		windowKey := location + "-" + date + "-" + value
		count := reservation.CountMap[windowKey]
		availability := maxOccupancy - count
		result[windowKey] = availability
	}
	return result
}
