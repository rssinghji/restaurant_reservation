package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testDefaultEndpoint(method string, endpoint string, request []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(request))
	writer := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleDefault)
	handler.ServeHTTP(writer, req)

	log.Printf("writer resp code [%d], writer resp [%s]", writer.Code, writer.Body.String())
	return writer
}

func TestDefaultEndpoint(test *testing.T) {
	writer := testDefaultEndpoint("GET", "/", nil)

	if http.StatusOK != writer.Code {
		test.Error("expected", http.StatusOK, "got", writer.Code)
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func testAddEndpoint(method string, endpoint string, request []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(request))
	writer := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleAddReservation)
	handler.ServeHTTP(writer, req)
	return writer
}

func TestAddReservationEndpoint(test *testing.T) {
	writer := testAddEndpoint("POST", "/add", []byte(`{
			"name": "Name 1",
			"date": "2024-08-18",
			"time": "11:05 AM",
			"email": "test@test.com"
		}`))

	if http.StatusOK != writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestAddReservationEndpointBad(test *testing.T) {
	writer := testAddEndpoint("GET", "/add", []byte(`{
			"name": "Name 1",
			"date": "2024-08-18",
			"time": "11:05 AM",
			"email": "test@test.com"
		}`))

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestAddReservationEndpointBad1(test *testing.T) {
	writer := testAddEndpoint("POST", "/add", []byte(`{
			"name": "Name 1",
			"date": "2024-08-18
			"time": "11:05 AM
			"email": "test@test.com"
		}`))

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestAddReservationEndpointBad2(test *testing.T) {
	writer := testAddEndpoint("POST", "/add", []byte(`{
			"name": "Name 1",
			"date": "2024-08-18,
			"time": "11:05 PM,
			"email": "test@test.com"
		}`))

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func testConfirmEndpoint(method string, endpoint string, request []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(request))
	writer := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleConfirmReservation)
	handler.ServeHTTP(writer, req)
	return writer
}

func TestConfirmReservationEndpoint(test *testing.T) {
	writer := testConfirmEndpoint("POST", "/confirm", []byte(`{
			"name": "Name 1",
			"date": "2024-08-18",
			"time": "11:05 AM",
			"email": "test@test.com"
		}`))

	if http.StatusOK != writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestConfirmReservationEndpointBad(test *testing.T) {
	writer := testConfirmEndpoint("GET", "/confirm", []byte(`{
			"name": "Name 1",
			"date": "2024-08-18",
			"time": "11:05 AM",
			"email": "test@test.com"
		}`))

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestConfirmReservationEndpointBad2(test *testing.T) {
	writer := testConfirmEndpoint("POST", "/confirm", []byte(`{
			name": "Name 9",
			"date": "2024-08-18",
			"time": "11:05 AM",
			"email": "test@test.com"
		}`))

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func testCancelEndpoint(method string, endpoint string, request []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(request))
	writer := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleCancelReservation)
	handler.ServeHTTP(writer, req)
	return writer
}

func TestCancelReservationEndpoint(test *testing.T) {
	writer := testCancelEndpoint("POST", "/cancel", []byte(`{
			"name": "Name 1",
			"date": "2024-08-18",
			"time": "11:05 AM",
			"email": "test@test.com"
		}`))

	if http.StatusOK != writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestCancelReservationEndpointBad(test *testing.T) {
	writer := testCancelEndpoint("GET", "/cancel", []byte(`{
			"name": "Name 1",
			"date": "2024-08-18",
			"time": "11:05 AM",
			"email": "test@test.com"
		}`))

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestCancelReservationEndpointBad1(test *testing.T) {
	writer := testCancelEndpoint("POST", "/cancel", []byte(`{
			"name": "Name 1",
			"date": "2024-08
			"time": "11:05 
			"email": "test@test.co
		}`))

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func testAvailabilityEndpoint(method string, endpoint string, request []byte, date string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(request))
	query := req.URL.Query()
	query.Add("location", "default")
	query.Add("date", date)
	req.URL.RawQuery = query.Encode()
	writer := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleShowAvailability)
	handler.ServeHTTP(writer, req)

	return writer
}

func TestShowAvailabilityEndpoint(test *testing.T) {
	writer := testAvailabilityEndpoint("GET", "/availability", nil, "2024-08-31")

	if http.StatusOK != writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestShowAvailabilityEndpointBad(test *testing.T) {
	writer := testAvailabilityEndpoint("GET", "/availability", nil, "vrevew-cwve-cwew")

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestShowAvailabilityEndpointBad1(test *testing.T) {
	writer := testAvailabilityEndpoint("POST", "/availability", nil, "2024-08-31")

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func testViewWaitingListEndpoint(method string, endpoint string, request []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(request))
	writer := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleViewWaitingList)
	handler.ServeHTTP(writer, req)

	return writer
}

func TestViewWaitingListEndpoint(test *testing.T) {
	writer := testViewWaitingListEndpoint("GET", "/waitinglist", nil)

	if http.StatusOK != writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestViewWaitingListEndpointBad(test *testing.T) {
	writer := testViewWaitingListEndpoint("POST", "/waitinglist", nil)

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func testViewReservationsEndpoint(method string, endpoint string, request []byte, date string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(request))
	query := req.URL.Query()
	query.Add("date", date)
	req.URL.RawQuery = query.Encode()
	writer := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleViewReservations)
	handler.ServeHTTP(writer, req)

	return writer
}

func TestViewReservationEndpoint(test *testing.T) {
	writer := testViewReservationsEndpoint("GET", "/view", nil, "2024-08-31")

	if http.StatusOK != writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestViewReservationEndpointBad(test *testing.T) {
	writer := testViewReservationsEndpoint("GET", "/view", nil, "csacdasd")

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}

func TestViewReservationEndpointBad1(test *testing.T) {
	writer := testViewReservationsEndpoint("POST", "/view", nil, "2024-08-31")

	if http.StatusOK == writer.Code {
		test.Fail()
	}
	body := writer.Body.String()
	if len(body) == 0 {
		test.Fail()
	}
}
