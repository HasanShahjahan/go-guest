package models

import "time"

type Guest struct {
	ID                 int       `json:"id"`
	Table              int       `json:"table"`
	Name               string    `json:"name"`
	AccompanyingGuests int       `json:"accompanying_guests"`
	Status             string    `json:"status"`
	ArrivalTime        time.Time `json:"time_arrived"`
}
