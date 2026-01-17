package controller

import (
	"HotelService/application/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ClientController struct {
	hotelService service.HotelService
}

func NewClientController(hotelService service.HotelService) *ClientController {
	return &ClientController{
		hotelService: hotelService,
	}
}

// ListHotels GET /client/hotels
func (c *ClientController) ListHotels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hotels, err := c.hotelService.ListHotels(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotels)
}

// GetHotelDetails GET /client/hotels/{id}
func (c *ClientController) GetHotelDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/client/hotels/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	hotel, err := c.hotelService.GetHotel(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotel)
}

// FindAvailableRooms GET /client/rooms/available
func (c *ClientController) FindAvailableRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rooms, err := c.hotelService.FindAvailableRooms(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

// test
func (c *ClientController) writeJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}
