package controller

import (
	"HotelService/application/dto"
	"HotelService/application/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type HotelierController struct {
	hotelService service.HotelService
}

func NewHotelierController(hotelService service.HotelService) *HotelierController {
	return &HotelierController{
		hotelService: hotelService,
	}
}

type CreateHotelRequest struct {
	Name    string              `json:"name"`
	Address string              `json:"address"`
	Rooms   []CreateRoomRequest `json:"rooms"`
}

type CreateRoomRequest struct {
	Number    string  `json:"number"`
	Type      string  `json:"type"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
}

// CreateHotel POST /hotelier/hotels
func (c *HotelierController) CreateHotel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateHotelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	rooms := make([]dto.RoomInput, len(req.Rooms))
	for i, room := range req.Rooms {
		rooms[i] = dto.RoomInput{
			Number:    room.Number,
			Type:      room.Type,
			Price:     room.Price,
			Available: room.Available,
		}
	}

	hotel, err := c.hotelService.CreateHotel(r.Context(), req.Name, req.Address, rooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(hotel)
}

type UpdateHotelRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// UpdateHotel PUT /hotelier/hotels/{id}
func (c *HotelierController) UpdateHotel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/hotelier/hotels/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var req UpdateHotelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	hotel, err := c.hotelService.UpdateHotel(r.Context(), id, req.Name, req.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotel)
}

// GetHotel GET /hotelier/hotels/{id}
func (c *HotelierController) GetHotel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/hotelier/hotels/")
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

type AddRoomRequest struct {
	Number    string  `json:"number"`
	Type      string  `json:"type"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
}

// AddRoom POST /hotelier/hotels/{hotelId}/rooms
func (c *HotelierController) AddRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/hotelier/hotels/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "rooms" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	hotelID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var req AddRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	room, err := c.hotelService.AddRoomToHotel(r.Context(), hotelID, req.Number, req.Type, req.Price, req.Available)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

type UpdateRoomRequest struct {
	Number    string  `json:"number"`
	Type      string  `json:"type"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
}

// UpdateRoom PUT /hotelier/rooms/{id}
func (c *HotelierController) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/hotelier/rooms/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	var req UpdateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	room, err := c.hotelService.UpdateRoom(r.Context(), id, req.Number, req.Type, req.Price, req.Available)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

// DeleteRoom DELETE /hotelier/rooms/{id}
func (c *HotelierController) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/hotelier/rooms/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	if err := c.hotelService.DeleteRoom(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateRoomAvailabilityRequest struct {
	Available bool `json:"available"`
}

// UpdateRoomAvailability PATCH /hotelier/rooms/{id}/availability
func (c *HotelierController) UpdateRoomAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/hotelier/rooms/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "availability" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	roomID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	var req UpdateRoomAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.hotelService.UpdateRoomAvailability(r.Context(), roomID, req.Available); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
