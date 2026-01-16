package usecase

import (
	"HotelService/domain/model"
	"context"
	"fmt"
	"time"
)

// refactor to service
type CreateHotelUseCase struct {
	hotelRepo HotelRepository
	roomRepo  RoomRepository
}

func NewCreateHotelUseCase(hotelRepo HotelRepository, roomRepo RoomRepository) *CreateHotelUseCase {
	return &CreateHotelUseCase{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
	}
}

func (uc *CreateHotelUseCase) Execute(ctx context.Context, name, address string, rooms []RoomInput) (*model.Hotel, error) {
	if name == "" {
		return nil, fmt.Errorf("hotel name is required")
	}
	if address == "" {
		return nil, fmt.Errorf("hotel address is required")
	}

	now := time.Now()

	hotel := &model.Hotel{
		Name:      name,
		Address:   address,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.hotelRepo.Save(ctx, hotel); err != nil {
		return nil, fmt.Errorf("failed to create hotel: %w", err)
	}

	if len(rooms) > 0 {
		for _, roomInput := range rooms {
			if roomInput.Number == "" {
				return nil, fmt.Errorf("room number is required")
			}
			if roomInput.Type == "" {
				return nil, fmt.Errorf("room type is required")
			}
			if roomInput.Price <= 0 {
				return nil, fmt.Errorf("room price must be positive")
			}

			room := &model.Room{
				HotelID:   hotel.ID,
				Number:    roomInput.Number,
				Type:      roomInput.Type,
				Price:     roomInput.Price,
				Available: roomInput.Available,
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := uc.roomRepo.Save(ctx, room); err != nil {
				return nil, fmt.Errorf("failed to create room %s: %w", roomInput.Number, err)
			}

			hotel.Rooms = append(hotel.Rooms, *room)
		}
	}

	return hotel, nil
}

type RoomInput struct {
	Number    string
	Type      string
	Price     float64
	Available bool
}

type GetHotelUseCase struct {
	hotelRepo HotelRepository
}

func NewGetHotelUseCase(hotelRepo HotelRepository) *GetHotelUseCase {
	return &GetHotelUseCase{hotelRepo: hotelRepo}
}

func (uc *GetHotelUseCase) Execute(ctx context.Context, id int64) (*model.Hotel, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid hotel ID")
	}

	hotel, err := uc.hotelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel: %w", err)
	}

	return hotel, nil
}

type ListHotelsUseCase struct {
	hotelRepo HotelRepository
}

func NewListHotelsUseCase(hotelRepo HotelRepository) *ListHotelsUseCase {
	return &ListHotelsUseCase{hotelRepo: hotelRepo}
}

func (uc *ListHotelsUseCase) Execute(ctx context.Context) ([]*model.Hotel, error) {
	hotels, err := uc.hotelRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list hotels: %w", err)
	}

	return hotels, nil
}

type UpdateRoomAvailabilityUseCase struct {
	roomRepo RoomRepository
}

func NewUpdateRoomAvailabilityUseCase(roomRepo RoomRepository) *UpdateRoomAvailabilityUseCase {
	return &UpdateRoomAvailabilityUseCase{roomRepo: roomRepo}
}

func (uc *UpdateRoomAvailabilityUseCase) Execute(ctx context.Context, roomID int64, available bool) error {
	if roomID <= 0 {
		return fmt.Errorf("invalid room ID")
	}

	_, err := uc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if err := uc.roomRepo.UpdateAvailability(ctx, roomID, available); err != nil {
		return fmt.Errorf("failed to update room availability: %w", err)
	}

	return nil
}

type FindAvailableRoomsUseCase struct {
	roomRepo RoomRepository
}

func NewFindAvailableRoomsUseCase(roomRepo RoomRepository) *FindAvailableRoomsUseCase {
	return &FindAvailableRoomsUseCase{roomRepo: roomRepo}
}

func (uc *FindAvailableRoomsUseCase) Execute(ctx context.Context) ([]*model.Room, error) {
	rooms, err := uc.roomRepo.FindAllAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find available rooms: %w", err)
	}

	return rooms, nil
}

// Hotelier Use Cases

type UpdateHotelUseCase struct {
	hotelRepo HotelRepository
}

func NewUpdateHotelUseCase(hotelRepo HotelRepository) *UpdateHotelUseCase {
	return &UpdateHotelUseCase{hotelRepo: hotelRepo}
}

func (uc *UpdateHotelUseCase) Execute(ctx context.Context, id int64, name, address string) (*model.Hotel, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid hotel ID")
	}
	if name == "" {
		return nil, fmt.Errorf("hotel name is required")
	}
	if address == "" {
		return nil, fmt.Errorf("hotel address is required")
	}

	existingHotel, err := uc.hotelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("hotel not found: %w", err)
	}

	existingHotel.Name = name
	existingHotel.Address = address
	existingHotel.UpdatedAt = time.Now()

	if err := uc.hotelRepo.Update(ctx, existingHotel); err != nil {
		return nil, fmt.Errorf("failed to update hotel: %w", err)
	}

	return existingHotel, nil
}

type AddRoomToHotelUseCase struct {
	roomRepo RoomRepository
}

func NewAddRoomToHotelUseCase(roomRepo RoomRepository) *AddRoomToHotelUseCase {
	return &AddRoomToHotelUseCase{roomRepo: roomRepo}
}

func (uc *AddRoomToHotelUseCase) Execute(ctx context.Context, hotelID int64, number, roomType string, price float64, available bool) (*model.Room, error) {
	if hotelID <= 0 {
		return nil, fmt.Errorf("invalid hotel ID")
	}
	if number == "" {
		return nil, fmt.Errorf("room number is required")
	}
	if roomType == "" {
		return nil, fmt.Errorf("room type is required")
	}
	if price <= 0 {
		return nil, fmt.Errorf("room price must be positive")
	}

	now := time.Now()
	room := &model.Room{
		HotelID:   hotelID,
		Number:    number,
		Type:      roomType,
		Price:     price,
		Available: available,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.roomRepo.Save(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to add room: %w", err)
	}

	return room, nil
}

type UpdateRoomUseCase struct {
	roomRepo RoomRepository
}

func NewUpdateRoomUseCase(roomRepo RoomRepository) *UpdateRoomUseCase {
	return &UpdateRoomUseCase{roomRepo: roomRepo}
}

func (uc *UpdateRoomUseCase) Execute(ctx context.Context, id int64, number, roomType string, price float64, available bool) (*model.Room, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid room ID")
	}
	if number == "" {
		return nil, fmt.Errorf("room number is required")
	}
	if roomType == "" {
		return nil, fmt.Errorf("room type is required")
	}
	if price <= 0 {
		return nil, fmt.Errorf("room price must be positive")
	}

	existingRoom, err := uc.roomRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	existingRoom.Number = number
	existingRoom.Type = roomType
	existingRoom.Price = price
	existingRoom.Available = available
	existingRoom.UpdatedAt = time.Now()

	if err := uc.roomRepo.Update(ctx, existingRoom); err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	return existingRoom, nil
}

type DeleteRoomUseCase struct {
	roomRepo RoomRepository
}

func NewDeleteRoomUseCase(roomRepo RoomRepository) *DeleteRoomUseCase {
	return &DeleteRoomUseCase{roomRepo: roomRepo}
}

func (uc *DeleteRoomUseCase) Execute(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid room ID")
	}

	_, err := uc.roomRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if err := uc.roomRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	return nil
}

// Client Use Cases

type GetHotelDetailsUseCase struct {
	hotelRepo HotelRepository
}

func NewGetHotelDetailsUseCase(hotelRepo HotelRepository) *GetHotelDetailsUseCase {
	return &GetHotelDetailsUseCase{hotelRepo: hotelRepo}
}

func (uc *GetHotelDetailsUseCase) Execute(ctx context.Context, id int64) (*model.Hotel, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid hotel ID")
	}

	hotel, err := uc.hotelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel details: %w", err)
	}

	return hotel, nil
}
