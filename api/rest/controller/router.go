package controller

import (
	"HotelService/application/usecase"
	"HotelService/infrastructure/db"
	"database/sql"
	"net/http"
	"strings"
)

func SetupRoutes(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	hotelRepo := db.NewHotelRepository(db)
	roomRepo := db.NewRoomRepository(db)

	createHotelUC := usecase.NewCreateHotelUseCase(hotelRepo, roomRepo)
	updateHotelUC := usecase.NewUpdateHotelUseCase(hotelRepo)
	getHotelUC := usecase.NewGetHotelUseCase(hotelRepo)
	listHotelsUC := usecase.NewListHotelsUseCase(hotelRepo)
	addRoomUC := usecase.NewAddRoomToHotelUseCase(roomRepo)
	updateRoomUC := usecase.NewUpdateRoomUseCase(roomRepo)
	deleteRoomUC := usecase.NewDeleteRoomUseCase(roomRepo)
	updateAvailabilityUC := usecase.NewUpdateRoomAvailabilityUseCase(roomRepo)
	findAvailableRoomsUC := usecase.NewFindAvailableRoomsUseCase(roomRepo)
	getHotelDetailsUC := usecase.NewGetHotelDetailsUseCase(hotelRepo)

	hotelierCtrl := NewHotelierController(
		createHotelUC,
		updateHotelUC,
		getHotelUC,
		addRoomUC,
		updateRoomUC,
		deleteRoomUC,
		updateAvailabilityUC,
	)
	clientCtrl := NewClientController(
		listHotelsUC,
		getHotelDetailsUC,
		findAvailableRoomsUC,
	)

	// Hotelier routes
	mux.HandleFunc("/hotelier/hotels", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			hotelierCtrl.CreateHotel(w, r)
		} else if r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/") {
			hotelierCtrl.GetHotel(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/hotelier/hotels/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			hotelierCtrl.UpdateHotel(w, r)
		} else if strings.Contains(r.URL.Path, "/rooms") {
			if r.Method == http.MethodPost {
				hotelierCtrl.AddRoom(w, r)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/hotelier/rooms/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			hotelierCtrl.UpdateRoom(w, r)
		} else if r.Method == http.MethodDelete {
			hotelierCtrl.DeleteRoom(w, r)
		} else if r.Method == http.MethodPatch && strings.Contains(r.URL.Path, "/availability") {
			hotelierCtrl.UpdateRoomAvailability(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Client routes
	mux.HandleFunc("/client/hotels", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if strings.Contains(r.URL.Path, "/") {
				clientCtrl.GetHotelDetails(w, r)
			} else {
				clientCtrl.ListHotels(w, r)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/client/hotels/", clientCtrl.GetHotelDetails)
	mux.HandleFunc("/client/rooms/available", clientCtrl.FindAvailableRooms)

	return mux
}
