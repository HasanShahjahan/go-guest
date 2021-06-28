package seed

import (
	"database/sql"
	"fmt"
	"log"
)

func Load(db *sql.DB) {
	fmt.Println("Loading seed data...")
	clearTable(db)
	ensureTableExists(db)
	setupSeedData(db)
	fmt.Println("Loading seed data is successful.")
}

func clearTable(db *sql.DB) {
	if _, err := db.Exec("DROP TABLE guest"); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec("DROP TABLE accommodation"); err != nil {
		log.Fatal(err)
	}
}

func ensureTableExists(db *sql.DB) {
	if _, err := db.Exec(accommodationTableCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(guestTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func setupSeedData(db *sql.DB) {
	if _, err := db.Exec(insertAccommodation); err != nil {
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
