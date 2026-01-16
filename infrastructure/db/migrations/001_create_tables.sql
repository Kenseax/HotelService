-- Create hotels table
CREATE TABLE IF NOT EXISTS hotels (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create rooms table
CREATE TABLE IF NOT EXISTS rooms (
    id BIGSERIAL PRIMARY KEY,
    hotel_id BIGINT NOT NULL,
    number VARCHAR(50) NOT NULL,
    type VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_hotel
        FOREIGN KEY (hotel_id)
        REFERENCES hotels(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_room_number_per_hotel
        UNIQUE (hotel_id, number)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_rooms_hotel_id ON rooms(hotel_id);
CREATE INDEX IF NOT EXISTS idx_rooms_available ON rooms(available);
CREATE INDEX IF NOT EXISTS idx_rooms_type ON rooms(type);
CREATE INDEX IF NOT EXISTS idx_rooms_price ON rooms(price);
CREATE INDEX IF NOT EXISTS idx_hotels_created_at ON hotels(created_at);

-- Insert sample data (optional, remove if not needed)
INSERT INTO hotels (name, address) VALUES 
    ('Grand Hotel', 'Sudino, Glavnaya ul, 12'),
    ('Ocean View Resort', '456 Beach Blvd, Miami, FL');

INSERT INTO rooms (hotel_id, number, type, price, available) VALUES
    (1, '101', 'Single', 100.00, true),
    (1, '102', 'Double', 150.00, true),
    (1, '103', 'Suite', 250.00, false),
    (2, '201', 'Single', 120.00, true),
    (2, '202', 'Double', 180.00, true);
