package service

import (
	"HotelService/application/dto"
	"HotelService/domain/model"
	"context"
	"fmt"
	"time"
)

type HotelServiceImpl struct {
	hotelRepo HotelRepository
	roomRepo  RoomRepository
}

func NewHotelService(hotelRepo HotelRepository, roomRepo RoomRepository) HotelService {
	return &HotelServiceImpl{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
	}
}

func (s *HotelServiceImpl) CreateHotel(ctx context.Context, name, address string, rooms []dto.RoomInput) (*model.Hotel, error) {
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

	if err := s.hotelRepo.Save(ctx, hotel); err != nil {
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

			if err := s.roomRepo.Save(ctx, room); err != nil {
				return nil, fmt.Errorf("failed to create room %s: %w", roomInput.Number, err)
			}

			hotel.Rooms = append(hotel.Rooms, *room)
		}
	}

	return hotel, nil
}

func (s *HotelServiceImpl) GetHotel(ctx context.Context, id int64) (*model.Hotel, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid hotel ID")
	}

	hotel, err := s.hotelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel: %w", err)
	}

	return hotel, nil
}

func (s *HotelServiceImpl) ListHotels(ctx context.Context) ([]*model.Hotel, error) {
	hotels, err := s.hotelRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list hotels: %w", err)
	}

	return hotels, nil
}

func (s *HotelServiceImpl) UpdateRoomAvailability(ctx context.Context, roomID int64, available bool) error {
	if roomID <= 0 {
		return fmt.Errorf("invalid room ID")
	}

	_, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if err := s.roomRepo.UpdateAvailability(ctx, roomID, available); err != nil {
		return fmt.Errorf("failed to update room availability: %w", err)
	}

	return nil
}

func (s *HotelServiceImpl) FindAvailableRooms(ctx context.Context) ([]*model.Room, error) {
	rooms, err := s.roomRepo.FindAllAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find available rooms: %w", err)
	}

	return rooms, nil
}

func (s *HotelServiceImpl) UpdateHotel(ctx context.Context, id int64, name, address string) (*model.Hotel, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid hotel ID")
	}
	if name == "" {
		return nil, fmt.Errorf("hotel name is required")
	}
	if address == "" {
		return nil, fmt.Errorf("hotel address is required")
	}

	existingHotel, err := s.hotelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("hotel not found: %w", err)
	}

	existingHotel.Name = name
	existingHotel.Address = address
	existingHotel.UpdatedAt = time.Now()

	if err := s.hotelRepo.Update(ctx, existingHotel); err != nil {
		return nil, fmt.Errorf("failed to update hotel: %w", err)
	}

	return existingHotel, nil
}

func (s *HotelServiceImpl) AddRoomToHotel(ctx context.Context, hotelID int64, number, roomType string, price float64, available bool) (*model.Room, error) {
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

	if err := s.roomRepo.Save(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to add room: %w", err)
	}

	return room, nil
}

func (s *HotelServiceImpl) UpdateRoom(ctx context.Context, id int64, number, roomType string, price float64, available bool) (*model.Room, error) {
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

	existingRoom, err := s.roomRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	existingRoom.Number = number
	existingRoom.Type = roomType
	existingRoom.Price = price
	existingRoom.Available = available
	existingRoom.UpdatedAt = time.Now()

	if err := s.roomRepo.Update(ctx, existingRoom); err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	return existingRoom, nil
}

func (s *HotelServiceImpl) DeleteRoom(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid room ID")
	}

	_, err := s.roomRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	if err := s.roomRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	return nil
}
