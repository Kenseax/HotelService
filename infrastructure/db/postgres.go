package db

import (
	"HotelService/domain/model"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type HotelPostgresRepository struct {
	db *sql.DB
}

type RoomPostgresRepository struct {
	db *sql.DB
}

func NewHotelRepository(db *sql.DB) *HotelPostgresRepository {
	return &HotelPostgresRepository{db: db}
}

func NewRoomRepository(db *sql.DB) *RoomPostgresRepository {
	return &RoomPostgresRepository{db: db}
}

// HotelPostgresRepository

func (r *HotelPostgresRepository) Save(ctx context.Context, hotel *model.Hotel) error {
	if hotel == nil {
		return fmt.Errorf("hotel cannot be nil")
	}

	now := time.Now()
	query := `
		INSERT INTO hotels (name, address, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		hotel.Name,
		hotel.Address,
		now,
		now,
	).Scan(&hotel.ID)

	if err != nil {
		return fmt.Errorf("failed to save hotel: %w", err)
	}

	hotel.CreatedAt = now
	hotel.UpdatedAt = now
	return nil
}

func (r *HotelPostgresRepository) Update(ctx context.Context, hotel *model.Hotel) error {
	if hotel == nil {
		return fmt.Errorf("hotel cannot be nil")
	}
	if hotel.ID == 0 {
		return fmt.Errorf("hotel ID is required for update")
	}

	query := `
		UPDATE hotels 
		SET name = $1, address = $2, updated_at = $3 
		WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		hotel.Name,
		hotel.Address,
		time.Now(),
		hotel.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update hotel: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("hotel with ID %d not found", hotel.ID)
	}

	hotel.UpdatedAt = time.Now()
	return nil
}

func (r *HotelPostgresRepository) FindByID(ctx context.Context, id int64) (*model.Hotel, error) {
	query := `
		SELECT id, name, address, created_at, updated_at 
		FROM hotels 
		WHERE id = $1`

	hotel := &model.Hotel{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&hotel.ID,
		&hotel.Name,
		&hotel.Address,
		&hotel.CreatedAt,
		&hotel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("hotel with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to find hotel: %w", err)
	}

	roomsQuery := `
		SELECT id, hotel_id, number, type, price, available, created_at, updated_at
		FROM rooms
		WHERE hotel_id = $1
		ORDER BY number`

	rows, err := r.db.QueryContext(ctx, roomsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load rooms: %w", err)
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Type, &room.Price, &room.Available, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rooms: %w", err)
	}

	hotel.Rooms = rooms
	return hotel, nil
}

func (r *HotelPostgresRepository) FindAll(ctx context.Context) ([]*model.Hotel, error) {
	query := `
		SELECT id, name, address, created_at, updated_at 
		FROM hotels 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all hotels: %w", err)
	}
	defer rows.Close()

	var hotels []*model.Hotel
	for rows.Next() {
		hotel := &model.Hotel{}
		if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Address, &hotel.CreatedAt, &hotel.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan hotel: %w", err)
		}
		hotels = append(hotels, hotel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating hotels: %w", err)
	}

	return hotels, nil
}

func (r *HotelPostgresRepository) Delete(ctx context.Context, id int64) error {

	query := `DELETE FROM hotels WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete hotel: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("hotel with ID %d not found", id)
	}

	return nil
}

// RoomPostgresRepository

func (r *RoomPostgresRepository) Save(ctx context.Context, room *model.Room) error {
	if room == nil {
		return fmt.Errorf("room cannot be nil")
	}

	query := `
		INSERT INTO rooms (hotel_id, number, type, price, available, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		room.HotelID,
		room.Number,
		room.Type,
		room.Price,
		room.Available,
		now,
		now,
	).Scan(&room.ID)

	if err != nil {
		return fmt.Errorf("failed to save room: %w", err)
	}

	room.CreatedAt = now
	room.UpdatedAt = now
	return nil
}

func (r *RoomPostgresRepository) Update(ctx context.Context, room *model.Room) error {
	if room == nil {
		return fmt.Errorf("room cannot be nil")
	}
	if room.ID == 0 {
		return fmt.Errorf("room ID is required for update")
	}

	query := `
		UPDATE rooms
		SET hotel_id = $1, number = $2, type = $3, price = $4, available = $5, updated_at = $6
		WHERE id = $7`

	result, err := r.db.ExecContext(ctx, query,
		room.HotelID,
		room.Number,
		room.Type,
		room.Price,
		room.Available,
		time.Now(),
		room.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("room with ID %d not found", room.ID)
	}

	room.UpdatedAt = time.Now()
	return nil
}

func (r *RoomPostgresRepository) FindByID(ctx context.Context, id int64) (*model.Room, error) {
	query := `
		SELECT id, hotel_id, number, type, price, available, created_at, updated_at
		FROM rooms
		WHERE id = $1`

	room := &model.Room{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.HotelID,
		&room.Number,
		&room.Type,
		&room.Price,
		&room.Available,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to find room: %w", err)
	}

	return room, nil
}

func (r *RoomPostgresRepository) FindAll(ctx context.Context) ([]*model.Room, error) {
	query := `
		SELECT id, hotel_id, number, type, price, available, created_at, updated_at
		FROM rooms
		ORDER BY hotel_id, number`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all rooms: %w", err)
	}
	defer rows.Close()

	return r.scanRooms(rows)
}

func (r *RoomPostgresRepository) FindByHotelID(ctx context.Context, hotelID int64) ([]*model.Room, error) {
	query := `
		SELECT id, hotel_id, number, type, price, available, created_at, updated_at
		FROM rooms
		WHERE hotel_id = $1
		ORDER BY number`

	rows, err := r.db.QueryContext(ctx, query, hotelID)
	if err != nil {
		return nil, fmt.Errorf("failed to find rooms by hotel ID: %w", err)
	}
	defer rows.Close()

	return r.scanRooms(rows)
}

func (r *RoomPostgresRepository) FindAllAvailable(ctx context.Context) ([]*model.Room, error) {
	query := `
		SELECT id, hotel_id, number, type, price, available, created_at, updated_at
		FROM rooms
		WHERE available = true
		ORDER BY hotel_id, number`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find available rooms: %w", err)
	}
	defer rows.Close()

	return r.scanRooms(rows)
}

func (r *RoomPostgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM rooms WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("room with ID %d not found", id)
	}

	return nil
}

func (r *RoomPostgresRepository) UpdateAvailability(ctx context.Context, id int64, available bool) error {
	query := `
		UPDATE rooms 
		SET available = $1, updated_at = $2 
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, available, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update room availability: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("room with ID %d not found", id)
	}

	return nil
}

// Repeatable room rows scan
func (r *RoomPostgresRepository) scanRooms(rows *sql.Rows) ([]*model.Room, error) {
	var rooms []*model.Room
	for rows.Next() {
		room := &model.Room{}
		if err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Type, &room.Price, &room.Available, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rooms: %w", err)
	}

	return rooms, nil
}
