package model

import "time"

type Hotel struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Rooms     []Room    `json:"rooms,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Room struct {
	ID        int64     `json:"id"`
	HotelID   int64     `json:"hotel_id"`
	Number    string    `json:"number"`
	Type      string    `json:"type"`
	Price     float64   `json:"price"`
	Available bool      `json:"available"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
