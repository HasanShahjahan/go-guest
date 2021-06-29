package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/HasanShahjahan/go-guest/api/responses"
	logging "github.com/HasanShahjahan/go-guest/api/utils"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	Upcoming = "upcoming"
	Attended = "attended"
	Archived = "archived"
)

type guest struct {
	ID                 int       `json:"id"`
	Table              int       `json:"table"`
	Name               string    `json:"name"`
	AccompanyingGuests int       `json:"accompanying_guests"`
	Status             string    `json:"status"`
	ArrivalTime        time.Time `json:"time_arrived"`
}

type accommodation struct {
	ID            int `json:"id"`
	TableNo       int `json:"table_no"`
	AvailableSeat int `json:"available_seat"`
	BookedSeat    int `json:"booked_seat"`
}

type seats struct {
	SeatsEmpty int `json:"seats_empty"`
}

type guestlist struct {
	guests []guest
}

func (server *Server) GetGuestLists(w http.ResponseWriter, r *http.Request) {
	guestLists, err := getGuestLists(server.DB)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		logging.Error(logTag, "Error during get guest lists", err)
		return
	}

	responses.RespondWithJSON(w, http.StatusOK, guestLists)
	logging.Info(logTag, "Get guest lists: ", guestLists)
}

func (server *Server) CreateGuest(w http.ResponseWriter, r *http.Request) {
	var g guest
	var ac accommodation
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		responses.RespondWithError(w, http.StatusBadRequest, "Invalid guest name")
		logging.Warn(logTag, "Invalid guest name")
		return
	}

	g = guest{Name: name, Status: Upcoming}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&g); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		logging.Error(logTag, "Error during decode request payload", err)
		return
	}
	defer r.Body.Close()

	logging.Info(logTag, "[Create Guest][Request]=%v", g)
	ac = accommodation{TableNo: g.Table}
	if err := ac.getAccommodationByTableNo(server.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			responses.RespondWithError(w, http.StatusNotFound, "Table is not found")
			logging.Error(logTag, "Invalid Table number, table=%d, error=%v", g.Table, err)
		default:
			responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
			logging.Error(logTag, "Error during get accommodation information, error=%v", err)
		}
		return
	}

	//Accompany guests and person itself
	if g.AccompanyingGuests+1 > ac.AvailableSeat-ac.BookedSeat {
		responses.RespondWithError(w, http.StatusUnprocessableEntity, "Insufficient space at the specified table")
		logging.Warn(logTag, "Insufficient space at the specified table, table=%d, error=%v", g.Table)
		return
	}

	g.Table = ac.ID
	ac.BookedSeat = ac.BookedSeat + g.AccompanyingGuests + 1
	if err := g.createGuest(server.DB, ac.BookedSeat); err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		logging.Error(logTag, "Guest creation is failed, error=%v", err)
		return
	}
	responses.RespondWithJSON(w, http.StatusCreated, g)
	logging.Info(logTag, "[Create Guest][Response]=%v", g)
}

func (server *Server) UpdateGuest(w http.ResponseWriter, r *http.Request) {
	var ac accommodation
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		responses.RespondWithError(w, http.StatusBadRequest, "Invalid guest name")
		logging.Warn(logTag, "Invalid guest name, name=%d", name)
		return
	}

	var g guest
	g = guest{Name: name}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&g); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		logging.Error(logTag, "Invalid request payload, error=%v", err)
		return
	}
	defer r.Body.Close()

	var databaseInfo guest
	databaseInfo = guest{Name: name}
	if err := databaseInfo.getGuest(server.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			responses.RespondWithError(w, http.StatusNotFound, "Guest is not exists.")
			logging.Error(logTag, "Guest is not exists, name=%d, error=%v", name, err)
		default:
			responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
			logging.Error(logTag, "Error during get guest, name=%d, error=%v", name, err)
		}
		return
	}

	ac = accommodation{TableNo: databaseInfo.Table}
	if err := ac.getAccommodationByTableId(server.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			responses.RespondWithError(w, http.StatusNotFound, "Table is not found")
			logging.Error(logTag, "Table is not found, table=%d, error=%v", ac.TableNo, err)
		default:
			responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
			logging.Error(logTag, "Error during get accommodation, table=%d, error=%v", ac.TableNo, err)
		}
		return
	}

	g.Table = databaseInfo.Table
	g.Status = Attended
	g.ArrivalTime = time.Now()

	if g.AccompanyingGuests > databaseInfo.AccompanyingGuests {
		ac.BookedSeat = ac.BookedSeat + (g.AccompanyingGuests - databaseInfo.AccompanyingGuests)
	} else if g.AccompanyingGuests < databaseInfo.AccompanyingGuests {
		ac.BookedSeat = ac.BookedSeat - (databaseInfo.AccompanyingGuests - g.AccompanyingGuests)
	}

	if ac.AvailableSeat < ac.BookedSeat {
		responses.RespondWithError(w, http.StatusUnprocessableEntity, "Insufficient space at the specified table")
		logging.Warn(logTag, "Insufficient space at the specified table")
		return
	}

	if err := g.updateGuest(server.DB, ac.BookedSeat); err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		logging.Error(logTag, "Guest update is failed ,error=%v", err)
		return
	}
	responses.RespondWithJSON(w, http.StatusOK, g)
}

func (server *Server) DeleteGuest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		responses.RespondWithError(w, http.StatusBadRequest, "Invalid guest name")
		logging.Error(logTag, "Invalid guest name ,name=%d", name)
		return
	}

	var g guest
	g = guest{Name: name}
	if err := g.getGuest(server.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			responses.RespondWithError(w, http.StatusNotFound, "Guest is not exists.")
			logging.Error(logTag, "Guest is not exists. ,name=%d", g.Name)
		default:
			responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
			logging.Error(logTag, "Error during get guest, error=%v", err)
		}
		return
	}

	var ac accommodation
	ac = accommodation{TableNo: g.Table}
	if err := ac.getAccommodationByTableId(server.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			responses.RespondWithError(w, http.StatusNotFound, "Table is not found")
			logging.Error(logTag, "Table is not found, table=%d", g.Table)
		default:
			responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
			logging.Error(logTag, "Error during get accommodation, tableNo=%d, error=%v", g.Table, err)
		}
		return
	}

	ac.BookedSeat = ac.BookedSeat - (g.AccompanyingGuests + 1)
	g.Status = Archived
	if err := g.deleteGuest(server.DB, ac.BookedSeat); err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		logging.Error(logTag, "Guest delete is failed ,error=%v", err)
		return
	}

	responses.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (server *Server) GetArrivedGuests(w http.ResponseWriter, r *http.Request) {
	var g guest
	g = guest{Status: Attended}
	guestLists, err := g.getGuests(server.DB)
	if err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		logging.Error(logTag, "Error during get arrived guests ,error=%v", err)
		return
	}

	responses.RespondWithJSON(w, http.StatusOK, guestLists)
}

func (server *Server) SeatsEmpty(w http.ResponseWriter, r *http.Request) {
	var s seats
	if err := s.seatCount(server.DB); err != nil {
		responses.RespondWithError(w, http.StatusInternalServerError, err.Error())
		logging.Error(logTag, "Error during get empty seats ,error=%v", err)
		return
	}
	responses.RespondWithJSON(w, http.StatusOK, s)
}

func getGuestLists(db *sql.DB) (*guestlist, error) {
	sqlStatement := `SELECT g.name ,a.table_no ,g.accompanying_guests 
                     FROM guest g 
                     JOIN accommodation a ON g.table_id = a.id `
	rows, err := db.Query(sqlStatement)

	if err != nil {
		logging.Error(logTag, "DB: Error during execute sql statement to get list  ,error=%v", err)
		return nil, err
	}
	defer rows.Close()

	var guests []guest
	for rows.Next() {
		var p guest
		if err := rows.Scan(&p.Name, &p.Table, &p.AccompanyingGuests); err != nil {
			logging.Error(logTag, "DB: Error during mapping of data  ,error=%v", err)
			return nil, err
		}
		guests = append(guests, p)
	}

	return &guestlist{guests: guests}, nil
}

func (a *guest) getGuests(db *sql.DB) (*guestlist, error) {
	rows, err := db.Query("SELECT g.name, g.accompanying_guests, g.arrival_time FROM guest g where g.status =?", a.Status)
	if err != nil {
		logging.Error(logTag, "DB: Error during sql statement to get arrived guest , error=%v", err)
		return nil, err
	}
	defer rows.Close()

	var guests []guest
	for rows.Next() {
		var p guest
		if err := rows.Scan(&p.Name, &p.AccompanyingGuests, &p.ArrivalTime); err != nil {
			logging.Error(logTag, "DB: Error during mapping of data  ,error=%v", err)
			return nil, err
		}
		guests = append(guests, p)
	}

	return &guestlist{guests: guests}, nil
}

func (a *guest) getGuest(db *sql.DB) error {
	return db.QueryRow("SELECT g.id , g.name , g.table_id, g.accompanying_guests, g.status FROM guest g WHERE g.name = ?", a.Name).Scan(&a.ID, &a.Name, &a.Table, &a.AccompanyingGuests, &a.Status)
}

func (a *guest) createGuest(db *sql.DB, bookedSeat int) error {
	ins, err := db.Prepare("INSERT INTO guest(name, table_id, accompanying_guests, status) VALUES(?, ?, ?, ?);")
	if err != nil {
		panic(err)
		logging.Error(logTag, "DB: Error during prepare create guest sql statement , error=%v", err)
	}
	defer ins.Close()

	res, err := ins.Exec(a.Name, a.Table, a.AccompanyingGuests, a.Status)
	rowsAffect, _ := res.RowsAffected()
	if err != nil || rowsAffect != 1 {
		logging.Error(logTag, "DB: Error inserting data, please check all fields. , error=%v", err)
		return err
	}

	update, err := db.Prepare("UPDATE accommodation SET booked_seat=? WHERE id=?;")
	if err != nil {
		panic(err)
		logging.Error(logTag, "DB: Error during update accommodation sql statement , error=%v", err)
	}
	defer update.Close()

	updaters, err := update.Exec(bookedSeat, a.Table)
	rowsAffected, _ := updaters.RowsAffected()
	if err != nil || rowsAffected == 0 {
		logging.Error(logTag, "DB: Error during update accommodation data., error=%v", err)
		return err
	}
	return nil
}

func (a *guest) updateGuest(db *sql.DB, bookedSeat int) error {
	update, err := db.Prepare(" UPDATE guest SET accompanying_guests=?, status=?, arrival_time=? WHERE name=?")
	if err != nil {
		panic(err)
		logging.Error(logTag, "DB: Error during update guest sql statement., error=%v", err)
	}
	defer update.Close()

	res, err := update.Exec(a.AccompanyingGuests, a.Status, a.ArrivalTime, a.Name)
	rowsAffect, _ := res.RowsAffected()
	if err != nil || rowsAffect != 1 {
		logging.Error(logTag, "DB: Error updating data, please check all fields., error=%v", err)
		return err
	}

	updateAcc, err := db.Prepare("UPDATE accommodation SET booked_seat=? WHERE id=?;")
	if err != nil {
		panic(err)
		logging.Error(logTag, "DB: Error updating accommodation data, please check all fields., error=%v", err)
	}
	defer updateAcc.Close()

	updaters, err := updateAcc.Exec(bookedSeat, a.Table)
	rowsAffected, _ := updaters.RowsAffected()
	if err != nil || rowsAffected == 0 {
		fmt.Printf("Error during update accommodation data.")
		return err
	}
	return nil
}

func (a *guest) deleteGuest(db *sql.DB, bookedSeat int) error {
	update, err := db.Prepare(" UPDATE guest SET status=? WHERE name=?")
	if err != nil {
		panic(err)
		logging.Error(logTag, "DB: Error updating delete sql preparation, please check all fields., error=%v", err)
	}
	defer update.Close()

	res, err := update.Exec(a.Status, a.Name)
	rowsAffect, _ := res.RowsAffected()
	if err != nil || rowsAffect != 1 {
		logging.Error(logTag, "DB: Error deleting data, please check all fields., error=%v", err)
		return err
	}

	updateAcc, err := db.Prepare("UPDATE accommodation SET booked_seat=? WHERE id=?;")
	if err != nil {
		panic(err)
		logging.Error(logTag, "DB: Error updating update accommodation sql preparation, please check all fields., error=%v", err)
	}
	defer updateAcc.Close()

	updaters, err := updateAcc.Exec(bookedSeat, a.Table)
	rowsAffected, _ := updaters.RowsAffected()
	if err != nil || rowsAffected == 0 {
		logging.Error(logTag, "DB: Error during update accommodation data., error=%v", err)
		return err
	}
	return nil
}

func (a *accommodation) getAccommodationByTableId(db *sql.DB) error {
	return db.QueryRow("SELECT a.id , a.available_seat , a.booked_seat FROM accommodation a WHERE a.id = ?", a.TableNo).Scan(&a.ID, &a.AvailableSeat, &a.BookedSeat)
}

func (a *accommodation) getAccommodationByTableNo(db *sql.DB) error {
	return db.QueryRow("SELECT a.id , a.available_seat , a.booked_seat FROM accommodation a WHERE a.table_no = ?", a.TableNo).Scan(&a.ID, &a.AvailableSeat, &a.BookedSeat)
}

func (s *seats) seatCount(db *sql.DB) error {
	return db.QueryRow("SELECT SUM(a.available_seat-a.booked_seat) as seats_empty FROM accommodation a").Scan(&s.SeatsEmpty)
}
