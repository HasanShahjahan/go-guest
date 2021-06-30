package controllers

import (
	"github.com/HasanShahjahan/go-guest/api/middlewares"
)

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/guest_list", server.GetGuestLists).Methods("GET")
	server.Router.HandleFunc("/guest_list/{name}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.CreateGuest))).Methods("POST")
	server.Router.HandleFunc("/guests/{name}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateGuest))).Methods("PUT")
	server.Router.HandleFunc("/guests/{name}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.DeleteGuest))).Methods("DELETE")
	server.Router.HandleFunc("/guests", server.GetArrivedGuests).Methods("GET")
	server.Router.HandleFunc("/seats_empty", server.SeatsEmpty).Methods("GET")
}
