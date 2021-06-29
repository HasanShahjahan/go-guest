package controllers

import (
	"database/sql"
	"github.com/HasanShahjahan/go-guest/api/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	logTag = "[Base]"
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
		logging.Warn(logTag, "Cannot connect to ", err)
	} else {
		logging.Info(logTag, "We are connected to the %s database \n", driver)
	}

	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	logging.Info(logTag, "Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
