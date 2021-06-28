package guest

import "net/http"

type guest interface {
	GetGuestLists(w http.ResponseWriter, r *http.Request)
	CreateGuest(w http.ResponseWriter, r *http.Request)
	UpdateGuest(w http.ResponseWriter, r *http.Request)
	DeleteGuest(w http.ResponseWriter, r *http.Request)
	GetArrivedGuests(w http.ResponseWriter, r *http.Request)
	SeatsEmpty(w http.ResponseWriter, r *http.Request)
}
