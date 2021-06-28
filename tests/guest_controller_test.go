package tests

import (
	"bytes"
	"encoding/json"
	"github.com/HasanShahjahan/go-guest/api/controllers"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a = controllers.Server{}

func TestEmptyTable(t *testing.T) {
	a.Initialize(os.Getenv("TEST_DB_DRIVER"), os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"))
	clearTable()
	ensureTableExists()

	req, _ := http.NewRequest("GET", "/guest_list", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "{}" {
		t.Errorf("Expected an empty object of array. Got %s", body)
	}
}

func TestCreateGuest(t *testing.T) {
	a.Initialize(os.Getenv("TEST_DB_DRIVER"), os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"))
	clearTable()
	ensureTableExists()
	setupSeedData()

	var jsonStr = []byte(`{"table": 1001,"accompanying_guests": 3}`)
	req, _ := http.NewRequest("POST", "/guest_list/Hasan", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Hasan" {
		t.Errorf("Expected guest name to be 'Hasan'. Got '%v'", m["name"])
	}
}

func TestUpdateGuest(t *testing.T) {
	a.Initialize(os.Getenv("TEST_DB_DRIVER"), os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"))
	clearTable()
	ensureTableExists()
	setupSeedData()
	addGuest()
	updateAccommodation()

	var jsonStr = []byte(`{"accompanying_guests": 6}`)
	req, _ := http.NewRequest("PUT", "/guests/Hasan", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Hasan" {
		t.Errorf("Expected guest name to be 'Hasan'. Got '%v'", m["name"])
	}
}

func TestDeleteGuest(t *testing.T) {
	a.Initialize(os.Getenv("TEST_DB_DRIVER"), os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"))
	clearTable()
	ensureTableExists()
	setupSeedData()
	addGuest()
	updateAccommodation()

	req, _ := http.NewRequest("DELETE", "/guests/Hasan", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["result"] != "success" {
		t.Errorf("Expected result to be 'success'. Got '%v'", m["result"])
	}
}

func TestArrivedGuests(t *testing.T) {
	a.Initialize(os.Getenv("TEST_DB_DRIVER"), os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"))
	clearTable()
	ensureTableExists()

	req, _ := http.NewRequest("GET", "/guests", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "{}" {
		t.Errorf("Expected an empty object of array. Got %s", body)
	}
}

func addGuest() {
	if _, err := a.DB.Exec(insertGuest); err != nil {
		log.Fatal(err)
	}
}

func updateAccommodation() {
	if _, err := a.DB.Exec(updateAcc); err != nil {
		log.Fatal(err)
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec(accommodationTableCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(guestTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	if _, err := a.DB.Exec("DROP TABLE guest"); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec("DROP TABLE accommodation"); err != nil {
		log.Fatal(err)
	}
}

func setupSeedData() {
	if _, err := a.DB.Exec(insertAccommodation); err != nil {
		log.Fatal(err)
	}
}

const guestTableCreationQuery = `CREATE TABLE IF NOT EXISTS guest
(
	id                     int(11)                                 NOT NULL     AUTO_INCREMENT,
	name                   varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
    table_id               int(11)                                 NOT NULL,
	accompanying_guests    int(11)                                 NOT NULL,
    status                 varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    arrival_time           datetime                                NULL,
	created_at             datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at             datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at             datetime                                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT guest_ibfk_1 FOREIGN KEY (table_id) REFERENCES accommodation (id),
	PRIMARY KEY (id)
)`

const accommodationTableCreationQuery = `CREATE TABLE IF NOT EXISTS accommodation 
(
	id                     int(11)    NOT NULL    AUTO_INCREMENT,
    table_no               int(11)    NOT NULL,
	available_seat         int(11)    NOT NULL,
	booked_seat            int(11)    NULL,
	PRIMARY KEY (id)
)`

const insertAccommodation = `INSERT INTO accommodation(table_no, available_seat, booked_seat) VALUES (1001,15,0), (1002,20,0), (1003,30,0);`

const insertGuest = `INSERT INTO guest(name, table_id, accompanying_guests, status) VALUES ('Hasan',1, 5, 'Upcoming');`

const updateAcc = `UPDATE accommodation SET booked_seat=6 WHERE table_no=1001;`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
