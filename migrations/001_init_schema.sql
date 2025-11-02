-- migrations/001_init_schema.sql

-- Table: yards
CREATE TABLE IF NOT EXISTS yards (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: blocks
CREATE TABLE IF NOT EXISTS blocks (
    id SERIAL PRIMARY KEY,
    yard_id INTEGER NOT NULL REFERENCES yards(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    max_slot INTEGER NOT NULL CHECK (max_slot > 0),
    max_row INTEGER NOT NULL CHECK (max_row > 0),
    max_tier INTEGER NOT NULL CHECK (max_tier > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(yard_id, code)
);

-- Table: yard_plans
-- Menyimpan perencanaan area untuk tipe container tertentu
CREATE TABLE IF NOT EXISTS yard_plans (
    id SERIAL PRIMARY KEY,
    block_id INTEGER NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
    slot_start INTEGER NOT NULL CHECK (slot_start > 0),
    slot_end INTEGER NOT NULL CHECK (slot_end >= slot_start),
    row_start INTEGER NOT NULL CHECK (row_start > 0),
    row_end INTEGER NOT NULL CHECK (row_end >= row_start),
    container_size INTEGER NOT NULL CHECK (container_size IN (20, 40)),
    container_height DECIMAL(3,1) NOT NULL CHECK (container_height IN (8.6, 9.6)),
    container_type VARCHAR(20) NOT NULL CHECK (container_type IN ('DRY', 'REEFER', 'OPEN_TOP')),
    stacking_priority VARCHAR(20) DEFAULT 'LEFT_TO_RIGHT',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk checking overlap
CREATE INDEX idx_yard_plans_block ON yard_plans(block_id);

-- Table: containers
-- Menyimpan container yang ada di yard
CREATE TABLE IF NOT EXISTS containers (
    id SERIAL PRIMARY KEY,
    container_number VARCHAR(50) UNIQUE NOT NULL,
    yard_id INTEGER NOT NULL REFERENCES yards(id),
    block_id INTEGER NOT NULL REFERENCES blocks(id),
    slot INTEGER NOT NULL,
    row INTEGER NOT NULL,
    tier INTEGER NOT NULL,
    container_size INTEGER NOT NULL CHECK (container_size IN (20, 40)),
    container_height DECIMAL(3,1) NOT NULL CHECK (container_height IN (8.6, 9.6)),
    container_type VARCHAR(20) NOT NULL CHECK (container_type IN ('DRY', 'REEFER', 'OPEN_TOP')),
    placed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(block_id, slot, row, tier),
    -- Untuk 40ft container, akan occupy 2 slots
    CHECK (
        (container_size = 20) OR 
        (container_size = 40 AND slot > 0)
    )
);

-- Index untuk performance
CREATE INDEX idx_containers_yard ON containers(yard_id);
CREATE INDEX idx_containers_block ON containers(block_id);
CREATE INDEX idx_containers_position ON containers(block_id, slot, row, tier);
CREATE INDEX idx_containers_number ON containers(container_number);

-- Seed data untuk testing
INSERT INTO yards (code, name, description) VALUES
('YRD1', 'Yard 1', 'Main container yard');

INSERT INTO blocks (yard_id, code, name, max_slot, max_row, max_tier) VALUES
(1, 'LC01', 'Loading Container Block 01', 10, 5, 5);

-- Yard Plans sesuai studi kasus
-- Plan untuk 20ft containers di slot 1-3, row 1-5
INSERT INTO yard_plans (block_id, slot_start, slot_end, row_start, row_end, container_size, container_height, container_type) VALUES
(1, 1, 3, 1, 5, 20, 8.6, 'DRY');

-- Plan untuk 40ft containers di slot 4-7, row 1-5
INSERT INTO yard_plans (block_id, slot_start, slot_end, row_start, row_end, container_size, container_height, container_type) VALUES
(1, 4, 7, 1, 5, 40, 8.6, 'DRY');