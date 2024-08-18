package api

import (
	util "reservation_system/utilities"
	"testing"
	"time"
)

func TestGetWindowBasic(test *testing.T) {
	window, err := getWindow("10:0 AM")
	if err != nil {
		test.Fail()
	}
	if window != "10AM-11AM" {
		test.Fail()
	}
}

func TestGetWindow(test *testing.T) {
	window, err := getWindow("11:05 AM")
	if err != nil {
		test.Fail()
	}
	if window != "11AM-12PM" {
		test.Fail()
	}
}

func TestGetWindowBad(test *testing.T) {
	_, err := getWindow("csdk : ks ")
	if err == nil {
		test.Fail()
	}
}

func TestGetWindow1(test *testing.T) {
	window, err := getWindow("12:00 PM")
	if err != nil {
		test.Fail()
	}
	if window != "12PM-1PM" {
		test.Fail()
	}
}

func TestGetWindow2(test *testing.T) {
	window, err := getWindow("1:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "1PM-2PM" {
		test.Fail()
	}
}

func TestGetWindow3(test *testing.T) {
	window, err := getWindow("2:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "2PM-3PM" {
		test.Fail()
	}
}

func TestGetWindow4(test *testing.T) {
	window, err := getWindow("3:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "3PM-4PM" {
		test.Fail()
	}
}

func TestGetWindow5(test *testing.T) {
	window, err := getWindow("3:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "3PM-4PM" {
		test.Fail()
	}
}

func TestGetWindow6(test *testing.T) {
	window, err := getWindow("4:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "4PM-5PM" {
		test.Fail()
	}
}

func TestGetWindow7(test *testing.T) {
	window, err := getWindow("5:59 PM")
	if err != nil {
		test.Fail()
	}
	if window != "5PM-6PM" {
		test.Fail()
	}
}

func TestGetWindow8(test *testing.T) {
	window, err := getWindow("6:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "6PM-7PM" {
		test.Fail()
	}
}

func TestGetWindow9(test *testing.T) {
	window, err := getWindow("7:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "7PM-8PM" {
		test.Fail()
	}
}

func TestGetWindow10(test *testing.T) {
	window, err := getWindow("8:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "8PM-9PM" {
		test.Fail()
	}
}

func TestGetWindow11(test *testing.T) {
	window, err := getWindow("9:05 PM")
	if err != nil {
		test.Fail()
	}
	if window != "9PM-10PM" {
		test.Fail()
	}
}

func TestGetWindow12(test *testing.T) {
	_, err := getWindow("10:0 PM")
	if err == nil {
		test.Fail()
	}
}

func TestGetmaxOccupancy(test *testing.T) {
	maxOccupancy := getMaxOccupancy("10AM-11AM")
	if maxOccupancy != 4 {
		test.Fail()
	}
}

func TestGetmaxOccupancyPeakTime(test *testing.T) {
	maxOccupancy := getMaxOccupancy("6PM-7PM")
	if maxOccupancy != 3 {
		test.Fail()
	}
}

func TestGetmaxOccupancyPeakTime1(test *testing.T) {
	maxOccupancy := getMaxOccupancy("7PM-8PM")
	if maxOccupancy != 3 {
		test.Fail()
	}
}

func TestIsValidRequestGood(test *testing.T) {
	isValid, err := IsRequestValid(util.ReservationDetails{Date: "2024-08-31", Time: "11:05 AM"})
	if err != nil || !isValid {
		test.Fail()
	}
}

func TestIsValidRequestBad(test *testing.T) {
	isValid, err := IsRequestValid(util.ReservationDetails{Date: "2024-08-11", Time: "11:05 AM"})
	if err == nil || isValid {
		test.Fail()
	}
}

func TestIsValidRequestBad1(test *testing.T) {
	isValid, err := IsRequestValid(util.ReservationDetails{Date: "2024-08-31", Time: "11:05 PM"})
	if err == nil || isValid {
		test.Fail()
	}
}

func TestIsValidRequestBad2(test *testing.T) {
	isValid, err := IsRequestValid(util.ReservationDetails{Date: "2024-08-31", Time: "9:05 AM"})
	if err == nil || isValid {
		test.Fail()
	}
}

func TestIsValidRequestBad3(test *testing.T) {
	isValid, err := IsRequestValid(util.ReservationDetails{Date: "cckwc", Time: "11:05 AM"})
	if err == nil || isValid {
		test.Fail()
	}
}

func TestGetAvailabilities(test *testing.T) {
	result := GetAvailabilities("default", "2024-08-18")
	if len(result) != 12 {
		test.Fail()
	}
	test.Log(len(result))
}

func TestViewReservationBYDates(test *testing.T) {
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"] = util.ReservationResponse{
		Name:        "test",
		Date:        "2024-08-31",
		Time:        "11:15 AM",
		Status:      "WAITING",
		Location:    "default",
		Email:       "test@gamil.com",
		RequestTime: time.Now(),
	}
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test2"] = util.ReservationResponse{
		Name:        "ztest2",
		Date:        "2024-08-31",
		Time:        "11:05 AM",
		Status:      "WAITING",
		Location:    "default",
		Email:       "test@gamil.com",
		RequestTime: time.Now(),
	}
	reservation.CountMap["default-2024-08-31-11AM-12PM"] = 2
	result, err := ViewReservationsByDate("default", "2024-08-31")
	if err != nil || len(result["11AM-12PM"]) != 2 {
		test.Fail()
	}

	// Cleanup keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestConfirmReservation(test *testing.T) {
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"] = util.ReservationResponse{
		Name:        "test",
		Date:        "2024-08-31",
		Time:        "11:15 AM",
		Status:      "WAITING",
		Location:    "default",
		Email:       "test@gamil.com",
		RequestTime: time.Now(),
	}
	reservation.CountMap["default-2024-08-31-11AM-12PM"] = 1
	confirmRequest := util.ConfirmReservation{Name: "test", Date: "2024-08-31", Time: "11:15 AM", Location: "default"}
	result, err := ConfirmReservation(confirmRequest)
	if err != nil || result != util.SuccessConfirmReservation {
		test.Fail()
	}
	if reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"].Status != "CONFIRMED" {
		test.Fail()
	}

	// Clean up keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestConfirmReservationBad(test *testing.T) {
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"] = util.ReservationResponse{
		Name:        "test",
		Date:        "2024-08-31",
		Time:        "11:15 AM",
		Status:      "WAITING",
		Location:    "default",
		Email:       "test@gamil.com",
		RequestTime: time.Now(),
	}
	reservation.CountMap["default-2024-08-31-11AM-12PM"] = 1
	confirmRequest := util.ConfirmReservation{Name: "test", Date: "2024-08-31", Time: "11:15 PM", Location: "default"}
	result, err := ConfirmReservation(confirmRequest)
	if err == nil || result == util.SuccessConfirmReservation {
		test.Fail()
	}

	// Clean up keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestConfirmReservationBadEmpty(test *testing.T) {
	confirmRequest := util.ConfirmReservation{Name: "test", Date: "2024-08-31", Time: "11:15 AM", Location: "default"}
	result, err := ConfirmReservation(confirmRequest)
	if err != nil || result != util.ErrorConfirmReservation {
		test.Fail()
	}

	// Clean up keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestCancelReservation(test *testing.T) {
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"] = util.ReservationResponse{
		Name:        "test",
		Date:        "2024-08-31",
		Time:        "11:15 AM",
		Status:      "WAITING",
		Location:    "default",
		Email:       "test@gamil.com",
		RequestTime: time.Now(),
	}
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test2"] = util.ReservationResponse{
		Name:        "ztest2",
		Date:        "2024-08-31",
		Time:        "11:05 AM",
		Status:      "WAITING",
		Location:    "default",
		Email:       "test@gamil.com",
		RequestTime: time.Now(),
	}
	reservation.CountMap["default-2024-08-31-11AM-12PM"] = 2
	request := util.ReservationDetails{Name: "test", Date: "2024-08-31", Time: "11:15 AM", Location: "default"}
	result, err := CancelReservation(request)
	if err != nil || result != util.SuccessCancellationMessage {
		test.Fail()
	}

	if reservation.CountMap["default-2024-08-31-11AM-12PM"] != 1 {
		test.Fail()
	}

	if _, ok := reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"]; ok {
		test.Fail()
	}

	// Cleanup keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestCancelReservationBad(test *testing.T) {
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"] = util.ReservationResponse{
		Name:        "test",
		Date:        "2024-08-31",
		Time:        "11:15 AM",
		Status:      "WAITING",
		Location:    "default",
		Email:       "test@gamil.com",
		RequestTime: time.Now(),
	}
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test2"] = util.ReservationResponse{
		Name:        "ztest2",
		Date:        "2024-08-31",
		Time:        "11:05 AM",
		Status:      "WAITING",
		Location:    "default",
		Email:       "test@gamil.com",
		RequestTime: time.Now(),
	}
	reservation.CountMap["default-2024-08-31-11AM-12PM"] = 2
	request := util.ReservationDetails{Name: "test", Date: "2024-08-31", Time: "11:15 PM", Location: "default"}
	result, err := CancelReservation(request)
	if err == nil || result == util.SuccessCancellationMessage {
		test.Fail()
	}

	if reservation.CountMap["default-2024-08-31-11AM-12PM"] == 1 {
		test.Fail()
	}

	if _, ok := reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"]; !ok {
		test.Fail()
	}

	// Cleanup keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestAddReservationBad(test *testing.T) {
	request := util.ReservationDetails{Name: "test", Date: "2024-08-31", Time: "11:15 PM", Location: "default", Email: "test@test.com"}
	_, err := AddReservation(request)
	if err == nil {
		test.Fail()
	}
	if _, ok := reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"]; ok {
		test.Fail()
	}
	if reservation.CountMap["default-2024-08-31-11AM-12PM"] == 1 {
		test.Fail()
	}

	// Cleanup keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestAddReservationMaxOccupancy(test *testing.T) {
	request := util.ReservationDetails{Name: "test", Date: "2024-08-31", Time: "11:15 AM", Location: "default", Email: "test@test.com"}
	reservation.CountMap["default-2024-08-31-11AM-12PM"] = 4
	_, err := AddReservation(request)
	if err != nil {
		test.Fail()
	}
	if _, ok := reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"]; ok {
		test.Fail()
	}
	if reservation.CountMap["default-2024-08-31-11AM-12PM"] != 4 {
		test.Fail()
	}

	// Cleanup keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestAutoCancellation(test *testing.T) {
	reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"] = util.ReservationResponse{
		Name:        "test",
		Date:        "2024-08-31",
		Time:        "11:15 AM",
		Status:      "WAITING",
		Location:    "default",
		RequestTime: time.Now(),
	}
	waitForConfirmation("default-2024-08-31-11AM-12PM", "default-2024-08-31-11AM-12PM-test", 10)

	time.Sleep(15 * time.Second)

	if _, ok := reservation.DetailsMap["default-2024-08-31-11AM-12PM-test"]; ok {
		test.Fail()
	}

	// Cleanup keys
	for k := range reservation.DetailsMap {
		delete(reservation.DetailsMap, k)
	}
	for k := range reservation.CountMap {
		delete(reservation.CountMap, k)
	}
}

func TestSendNotification(test *testing.T) {
	request := util.ReservationDetails{Name: "test", Date: "2024-08-31", Time: "11:15 AM", Location: "default"}
	err := sendNotification(request)
	if err != nil {
		test.Fail()
	}
}
