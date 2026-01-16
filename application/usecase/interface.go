package usecase

import (
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
