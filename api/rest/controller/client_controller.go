package controller

import (
	"HotelService/application/usecase"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ClientController struct {
	listHotelsUseCase         *usecase.ListHotelsUseCase
	getHotelDetailsUseCase    *usecase.GetHotelDetailsUseCase
	findAvailableRoomsUseCase *usecase.FindAvailableRoomsUseCase
}

func NewClientController(
	listHotelsUC *usecase.ListHotelsUseCase,
	getHotelDetailsUC *usecase.GetHotelDetailsUseCase,
	findAvailableRoomsUC *usecase.FindAvailableRoomsUseCase,
) *ClientController {
	return &ClientController{
		listHotelsUseCase:         listHotelsUC,
		getHotelDetailsUseCase:    getHotelDetailsUC,
		findAvailableRoomsUseCase: findAvailableRoomsUC,
	}
}

// ListHotels GET /client/hotels
func (c *ClientController) ListHotels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hotels, err := c.listHotelsUseCase.Execute(r.Context())
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

	hotel, err := c.getHotelDetailsUseCase.Execute(r.Context(), id)
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

	rooms, err := c.findAvailableRoomsUseCase.Execute(r.Context())
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
