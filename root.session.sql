-- Create the resource_types table with GORM UUID format
DROP TABLE IF EXISTS resource_types;
DROP TABLE IF EXISTS tnt_resource_metadata;
DROP TABLE IF EXISTS tnt_resource;

CREATE TABLE IF NOT EXISTS resource_types (
    resource_type_id VARCHAR(36) PRIMARY KEY,
    service_id VARCHAR(36) NOT NULL,
    name VARCHAR(45) NOT NULL,
    row_status INT DEFAULT 1,
    created_by VARCHAR(45),
    updated_by VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data
INSERT INTO resource_types (
    resource_type_id,
    service_id,
    name,
    row_status,
    created_by,
    updated_by,
    created_at,
    updated_at
) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    'a1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'User',
    1,
    'system_admin',
    'system_admin',
    '2024-01-01 10:00:00',
    '2024-01-01 10:00:00'
),
(
    '550e8400-e29b-41d4-a716-446655440001',
    'a1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'Group',
    1,
    'system_admin',
    'system_admin',
    '2024-01-01 10:30:00',
    '2024-01-01 10:30:00'
),
(
    '550e8400-e29b-41d4-a716-446655440002',
    'a1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'Tenant',
    1,
    'system_admin',
    'system_admin',
    '2024-01-01 11:00:00',
    '2024-01-01 11:00:00'
);

-- Create indexes for better performance
CREATE INDEX idx_resource_types_service_id ON resource_types(service_id);
CREATE INDEX idx_resource_types_name ON resource_types(name);