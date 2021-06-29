package controllers

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	Router *mux.Router
	DB     *sql.DB
}

func (server *Server) Initialize(driver, user, password, dbname string) {
	var err error
	var dataSourceName string
	if user == "" && password == "" && dbname == "" {
		dataSourceName = "root:password@/guests"
	} else {
		dataSourceName = user + ":" + password + "@/" + dbname
	}
	server.DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", driver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", "mysql")
	}

	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
