package api

import (
	"fmt"
	"log"
	"net/smtp"
	util "reservation_system/utilities"
	"strconv"
	"strings"
	"time"
)

var clockTime = map[string]int{"1": 13, "2": 14, "3": 15, "4": 16, "5": 17, "6": 18, "7": 19, "8": 20, "9": 21}

// manageWaitlist() : Goroutine to pull a waiting request for the canceled window on a FIFO basis
func manageWaitlist(window string) {
	log.Println("Starting waitlist routine")

	queue.Mutex.Lock()
	defer queue.Mutex.Unlock()
	for index, request := range queue.Queue {
		windowRequested := queue.Windows[index]
		if windowRequested != window {
			continue
		}
		_, err := AddReservation(request)
		if err != nil {
			log.Println("Problem in addinf the queued request", err)
		}

		// remove the request from the queue
		queue.Queue = append(queue.Queue[:index], queue.Queue[index+1:]...)
		queue.Windows = append(queue.Windows[:index], queue.Windows[index+1:]...)
		err = sendNotification(request)
		if err != nil {
			log.Println("Unable to send email notification to the requestor: ", err)
		}
		break
	}

	log.Println("No interested request found")
}

// sendNotification() : Example funciton to emulate sending notification
func sendNotification(req util.ReservationDetails) error {
	// send email notification to user
	log.Println("Sending email to the requestor at: ", req.Email)
	from := "testrestaurant@gmail.com"
	password := "some password"

	to := []string{req.Email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	emailMessage := fmt.Sprintf("Subject: Slot available \n\n Hi, %s,\n We've open slot for Location: %s, Date%s, Time%s. PLease confirm.",
		req.Name, req.Location, req.Date, req.Time)
	message := []byte(emailMessage)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Println("Ignoring the error as not using actual email")
	}
	return nil
}

// waitListCleanUp() : Goroutine to check every 30 minutes if some requests are invalid and removes the reequest
func waitListCleanUp(seconds int) {
	log.Println("Clearing up waitlist to start a fresh day")
	for {
		queue.Mutex.Lock()
		for index, value := range queue.Queue {
			// Check for date
			today := time.Now()
			requestedDate, _ := time.Parse("2006-01-02", value.Date)
			if today.After(requestedDate) {
				// Can't serve this request, remove
				queue.Queue = append(queue.Queue[:index], queue.Queue[index+1:]...)
				queue.Windows = append(queue.Windows[:index], queue.Windows[index+1:]...)
				continue
			}

			// Check for time
			requestedTime := value.Time
			requestSplits := strings.Split(requestedTime, ":")
			requestHour := requestSplits[0]
			hour := clockTime[requestHour]
			requestMinute := strings.Split(requestSplits[1], " ")[0]
			minute, _ := strconv.Atoi(requestMinute)
			currentHour := time.Now().Hour()
			currentMinutes := time.Now().Minute()
			if currentHour > hour || currentMinutes > minute || minute-currentMinutes < 0 || minute-currentMinutes < 10 {
				// Can't serve this request, remove
				queue.Queue = append(queue.Queue[:index], queue.Queue[index+1:]...)
				queue.Windows = append(queue.Windows[:index], queue.Windows[index+1:]...)
			}
		}
		queue.Mutex.Unlock()
		time.Sleep(time.Duration(time.Duration(seconds).Seconds())) // Check after every 30 minutes for cleanup
	}
}
