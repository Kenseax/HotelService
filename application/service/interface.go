package service

import (
	"HotelService/application/dto"
	"HotelService/domain/model"
	"context"
)

type HotelRepository interface {
	Save(ctx context.Context, hotel *model.Hotel) error
	Update(ctx context.Context, hotel *model.Hotel) error
	FindByID(ctx context.Context, id int64) (*model.Hotel, error)
	FindAll(ctx context.Context) ([]*model.Hotel, error)
	Delete(ctx context.Context, id int64) error
}

type RoomRepository interface {
	Save(ctx context.Context, room *model.Room) error
	Update(ctx context.Context, room *model.Room) error
	FindByID(ctx context.Context, id int64) (*model.Room, error)
	FindAll(ctx context.Context) ([]*model.Room, error)
	FindByHotelID(ctx context.Context, hotelID int64) ([]*model.Room, error)
	FindAllAvailable(ctx context.Context) ([]*model.Room, error)
	Delete(ctx context.Context, id int64) error
	UpdateAvailability(ctx context.Context, id int64, available bool) error
}

type HotelService interface {
	CreateHotel(ctx context.Context, name, address string, rooms []dto.RoomInput) (*model.Hotel, error)
	GetHotel(ctx context.Context, id int64) (*model.Hotel, error)
	ListHotels(ctx context.Context) ([]*model.Hotel, error)
	UpdateRoomAvailability(ctx context.Context, roomID int64, available bool) error
	FindAvailableRooms(ctx context.Context) ([]*model.Room, error)
	UpdateHotel(ctx context.Context, id int64, name, address string) (*model.Hotel, error)
	AddRoomToHotel(ctx context.Context, hotelID int64, number, roomType string, price float64, available bool) (*model.Room, error)
	UpdateRoom(ctx context.Context, id int64, number, roomType string, price float64, available bool) (*model.Room, error)
	DeleteRoom(ctx context.Context, id int64) error
}
