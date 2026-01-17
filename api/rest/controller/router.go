package controller

import (
	"HotelService/application/service"
	"HotelService/infrastructure/db"
	"database/sql"
	"net/http"
	"strings"
)

func SetupRoutes(conn *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	hotelRepo := db.NewHotelRepository(conn)
	roomRepo := db.NewRoomRepository(conn)

	hotelService := service.NewHotelService(hotelRepo, roomRepo)

	hotelierCtrl := NewHotelierController(hotelService)
	clientCtrl := NewClientController(hotelService)

	// Hotelier routes
	mux.HandleFunc("/hotelier/hotels", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			hotelierCtrl.CreateHotel(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/hotelier/hotels/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			hotelierCtrl.GetHotel(w, r)
		} else if r.Method == http.MethodPut {
			hotelierCtrl.UpdateHotel(w, r)
		} else if strings.Contains(r.URL.Path, "/rooms") && r.Method == http.MethodPost {
			hotelierCtrl.AddRoom(w, r)
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
	mux.HandleFunc("/client/hotels", clientCtrl.ListHotels)
	mux.HandleFunc("/client/hotels/", clientCtrl.GetHotelDetails)
	mux.HandleFunc("/client/rooms/available", clientCtrl.FindAvailableRooms)

	return mux
}
