package dtos

type GuestDto struct {
	Name string `json:"name"`
}

type guest struct {
	Table              int    `json:"table"`
	Name               string `json:"name"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

type Guestlist struct {
	Guests []guest `json:"guests"`
}
