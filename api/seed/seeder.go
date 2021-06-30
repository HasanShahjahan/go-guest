package seed

import (
	"database/sql"
	logging "github.com/HasanShahjahan/go-guest/api/utils"
	"log"
)

const (
	logTag = "[Seed]"
)

func Load(db *sql.DB) {
	logging.Info(logTag, "Loading seed data...")
	createDatabase(db)
	clearTable(db)
	ensureTableExists(db)
	setupSeedData(db)

}

func createDatabase(db *sql.DB) {
	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS guests;"); err != nil {
		logging.Error(logTag, "Error during database creation", err)
	}
	logging.Info(logTag, "Database creation is done.")
}
func clearTable(db *sql.DB) {
	if _, err := db.Exec("DROP TABLE IF EXISTS guests.guest"); err != nil {
		logging.Error(logTag, "Error during Guest table drop", err)
	}
	logging.Info(logTag, "Guest table is dropped.")

	if _, err := db.Exec("DROP TABLE IF EXISTS guests.accommodation"); err != nil {
		logging.Error(logTag, "Error during Accomodation table drop", err)
	}
	logging.Info(logTag, "Accommodation table is dropped. ")
}

func ensureTableExists(db *sql.DB) {
	if _, err := db.Exec(accommodationTableCreationQuery); err != nil {
		log.Fatal(err)
		logging.Error(logTag, "Error during Accomodation table creation", err)
	}
	logging.Info(logTag, "Accomodation table is created.")

	if _, err := db.Exec(guestTableCreationQuery); err != nil {
		log.Fatal(err)
		logging.Error(logTag, "Error during Accomodation table creation", err)
	}
	logging.Info(logTag, "Guest table is created.")
}

func setupSeedData(db *sql.DB) {
	if _, err := db.Exec(insertAccommodation); err != nil {
		log.Fatal(err)
	}
	logging.Info(logTag, "Seeding data is successful.")
}

const guestTableCreationQuery = `CREATE TABLE IF NOT EXISTS guests.guest
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

const accommodationTableCreationQuery = `CREATE TABLE IF NOT EXISTS guests.accommodation 
(
	id                     int(11)    NOT NULL    AUTO_INCREMENT,
    table_no               int(11)    NOT NULL,
	available_seat         int(11)    NOT NULL,
	booked_seat            int(11)    NULL,
	PRIMARY KEY (id)
)`

const insertAccommodation = `INSERT INTO guests.accommodation(table_no, available_seat, booked_seat) VALUES (1001,15,0), (1002,20,0), (1003,30,0);`
