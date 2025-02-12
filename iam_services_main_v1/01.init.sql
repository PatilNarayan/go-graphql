-- Create the mst_resource_types table with MySQL UUID format
DROP TABLE IF EXISTS mst_resource_types;

CREATE TABLE IF NOT EXISTS mst_resource_types (
    resource_type_id VARCHAR(36) PRIMARY KEY,
    service_id VARCHAR(36) NOT NULL,
    name VARCHAR(36) NOT NULL,
    row_status INT DEFAULT 1,
    created_by VARCHAR(36) NOT NULL,
    updated_by VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data with proper UUID values
INSERT INTO mst_resource_types (
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
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '2024-01-01 10:00:00',
    '2024-01-01 10:00:00'
),
(
    '550e8400-e29b-41d4-a716-446655440001',
    'b1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'Group',
    1,
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '2024-01-01 10:30:00',
    '2024-01-01 10:30:00'
),
(
    '550e8400-e29b-41d4-a716-446655440002',
    'c1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'Tenant',
    1,
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '2024-01-01 11:00:00',
    '2024-01-01 11:00:00'
),
(
    '550e8400-e29b-41d4-a716-446655440003',
    'd1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'Role',
    1,
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '2024-01-01 11:30:00',
    '2024-01-01 11:30:00'
),
(
    '550e8400-e29b-41d4-a716-446655440004',
    'e1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'Root',
    1,
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '2024-01-01 12:00:00',
    '2024-01-01 12:00:00'
),
(
    '550e8400-e29b-41d4-a716-446655440005',
    'f1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'Account',
    1,
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '2024-01-01 12:30:00',
    '2024-01-01 12:30:00'
),
(
    '550e8400-e29b-41d4-a716-446655440006',
    'g1b2c3d4-e5f6-4747-8899-aabbccddeeff',
    'Client Organization Unit',
    1,
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    '2024-01-01 13:00:00',
    '2024-01-01 13:00:00'
);

