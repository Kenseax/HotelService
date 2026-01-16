package main

import (
	"HotelService/application/usecase"
	"HotelService/infrastructure/db"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// 1. Setup database connection
	dbConfig := db.DefaultConfig()

	database, err := db.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	fmt.Println("✓ Connected to database successfully")

	// 2. Initialize repositories
	hotelRepo := db.NewHotelRepository(database)
	roomRepo := db.NewRoomRepository(database)

	fmt.Println("✓ Repositories initialized")

	// 3. Initialize use cases
	createHotelUC := usecase.NewCreateHotelUseCase(hotelRepo, roomRepo)
	updateHotelUC := usecase.NewUpdateHotelUseCase(hotelRepo)
	getHotelUC := usecase.NewGetHotelUseCase(hotelRepo)
	listHotelsUC := usecase.NewListHotelsUseCase(hotelRepo)
	addRoomUC := usecase.NewAddRoomToHotelUseCase(roomRepo)
	updateRoomUC := usecase.NewUpdateRoomUseCase(roomRepo)
	updateRoomAvailabilityUC := usecase.NewUpdateRoomAvailabilityUseCase(roomRepo)
	findAvailableRoomsUC := usecase.NewFindAvailableRoomsUseCase(roomRepo)

	fmt.Println("✓ Use cases initialized")

	// 4. Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5. Example: Create a new hotel with rooms
	fmt.Println("\n--- Creating a new hotel ---")
	rooms := []usecase.RoomInput{
		{Number: "301", Type: "Suite", Price: 300.00, Available: true},
		{Number: "302", Type: "Double", Price: 200.00, Available: true},
		{Number: "303", Type: "Single", Price: 150.00, Available: false},
	}
	hotel, err := createHotelUC.Execute(ctx, "Luxury Hotel", "789 Park Ave, Seattle, WA", rooms)
	if err != nil {
		log.Printf("Error creating hotel: %v", err)
	} else {
		fmt.Printf("✓ Created hotel: ID=%d, Name=%s, Rooms=%d\n", hotel.ID, hotel.Name, len(hotel.Rooms))
		for _, room := range hotel.Rooms {
			fmt.Printf("  - Room %s: %s ($%.2f, Available: %v)\n",
				room.Number, room.Type, room.Price, room.Available)
		}
	}

	// 6. Example: Get hotel by ID
	fmt.Println("\n--- Getting hotel by ID ---")
	if hotel != nil {
		fetchedHotel, err := getHotelUC.Execute(ctx, hotel.ID)
		if err != nil {
			log.Printf("Error getting hotel: %v", err)
		} else {
			fmt.Printf("✓ Retrieved hotel: ID=%d, Name=%s, Address=%s\n",
				fetchedHotel.ID, fetchedHotel.Name, fetchedHotel.Address)
			fmt.Printf("  Rooms: %d\n", len(fetchedHotel.Rooms))
			for _, room := range fetchedHotel.Rooms {
				fmt.Printf("    - Room %s (Available: %v)\n", room.Number, room.Available)
			}
		}
	}

	// 7. Example: List all hotels
	fmt.Println("\n--- Listing all hotels ---")
	hotels, err := listHotelsUC.Execute(ctx)
	if err != nil {
		log.Printf("Error listing hotels: %v", err)
	} else {
		fmt.Printf("✓ Found %d hotels:\n", len(hotels))
		for _, h := range hotels {
			fmt.Printf("  - ID=%d, Name=%s\n", h.ID, h.Name)
		}
	}

	// 8. Example: Find available rooms
	fmt.Println("\n--- Finding available rooms ---")
	availableRooms, err := findAvailableRoomsUC.Execute(ctx)
	if err != nil {
		log.Printf("Error finding available rooms: %v", err)
	} else {
		fmt.Printf("✓ Found %d available rooms:\n", len(availableRooms))
		for _, room := range availableRooms {
			fmt.Printf("  - Room %s (Hotel ID: %d)\n", room.Number, room.HotelID)
		}
	}

	// 9. Example: Update room availability
	fmt.Println("\n--- Updating room availability ---")
	if len(availableRooms) > 0 {
		roomID := availableRooms[0].ID
		err = updateRoomAvailabilityUC.Execute(ctx, roomID, false)
		if err != nil {
			log.Printf("Error updating room availability: %v", err)
		} else {
			fmt.Printf("✓ Room %d marked as unavailable\n", roomID)
		}

		// Mark it back as available
		err = updateRoomAvailabilityUC.Execute(ctx, roomID, true)
		if err != nil {
			log.Printf("Error updating room availability: %v", err)
		} else {
			fmt.Printf("✓ Room %d marked as available again\n", roomID)
		}
	}

	// 10. Example: Update hotel information (Hotelier operation)
	fmt.Println("\n--- Updating hotel information ---")
	if hotel != nil {
		updatedHotel, err := updateHotelUC.Execute(ctx, hotel.ID, "Updated Luxury Hotel", "999 Updated Ave, Seattle, WA")
		if err != nil {
			log.Printf("Error updating hotel: %v", err)
		} else {
			fmt.Printf("✓ Updated hotel: ID=%d, Name=%s, Address=%s\n",
				updatedHotel.ID, updatedHotel.Name, updatedHotel.Address)
		}
	}

	// 11. Example: Add a new room to hotel (Hotelier operation)
	fmt.Println("\n--- Adding a new room to hotel ---")
	if hotel != nil {
		newRoom, err := addRoomUC.Execute(ctx, hotel.ID, "304", "Deluxe", 250.00, true)
		if err != nil {
			log.Printf("Error adding room: %v", err)
		} else {
			fmt.Printf("✓ Added room: ID=%d, Number=%s, Type=%s, Price=$%.2f\n",
				newRoom.ID, newRoom.Number, newRoom.Type, newRoom.Price)
		}
	}

	// 12. Example: Update room information (Hotelier operation)
	fmt.Println("\n--- Updating room information ---")
	if hotel != nil {
		// Get the hotel again to see the new room
		updatedHotel, err := getHotelUC.Execute(ctx, hotel.ID)
		if err == nil && len(updatedHotel.Rooms) > 0 {
			roomToUpdate := updatedHotel.Rooms[len(updatedHotel.Rooms)-1] // Last room (newly added)
			updatedRoom, err := updateRoomUC.Execute(ctx, roomToUpdate.ID, "304", "Premium Suite", 350.00, true)
			if err != nil {
				log.Printf("Error updating room: %v", err)
			} else {
				fmt.Printf("✓ Updated room: ID=%d, Number=%s, Type=%s, Price=$%.2f\n",
					updatedRoom.ID, updatedRoom.Number, updatedRoom.Type, updatedRoom.Price)
			}
		}
	}

	fmt.Println("\n✓ All examples completed successfully!")
	fmt.Println("\n--- API Endpoints Summary ---")
	fmt.Println("Hotelier Endpoints:")
	fmt.Println("  POST   /hotelier/hotels                    - Create hotel")
	fmt.Println("  PUT    /hotelier/hotels/{id}               - Update hotel")
	fmt.Println("  GET    /hotelier/hotels/{id}               - Get hotel details")
	fmt.Println("  POST   /hotelier/hotels/{id}/rooms         - Add room to hotel")
	fmt.Println("  PUT    /hotelier/rooms/{id}                - Update room")
	fmt.Println("  DELETE /hotelier/rooms/{id}                - Delete room")
	fmt.Println("  PATCH  /hotelier/rooms/{id}/availability   - Update room availability")
	fmt.Println("\nClient Endpoints:")
	fmt.Println("  GET    /client/hotels                      - List all hotels")
	fmt.Println("  GET    /client/hotels/{id}                 - Get hotel details")
	fmt.Println("  GET    /client/rooms/available             - Find available rooms")
}
