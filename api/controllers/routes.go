package controllers

import logging "github.com/HasanShahjahan/go-guest/api/utils"

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/guest_list", server.GetGuestLists).Methods("GET")
	server.Router.HandleFunc("/guest_list/{name}", server.CreateGuest).Methods("POST")
	server.Router.HandleFunc("/guests/{name}", server.UpdateGuest).Methods("PUT")
	server.Router.HandleFunc("/guests/{name}", server.DeleteGuest).Methods("DELETE")
	server.Router.HandleFunc("/guests", server.GetArrivedGuests).Methods("GET")
	server.Router.HandleFunc("/seats_empty", server.SeatsEmpty).Methods("GET")

	logging.Info(logTag, "Routes is initialized successfully.")
}
