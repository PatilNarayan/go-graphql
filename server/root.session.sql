-- Insert dummy data for 
INSERT INTO tnt_resource (
    resource_id,
    parent_resource_id,
    resource_type_id,
    name,
    row_status,
    created_by,
    updated_by,
    created_at,
    updated_at
) VALUES
(
    '550e8400-1111-41d4-a716-446655441111', -- ResourceID
    '11111111-1111-1111-1111-111111111111', -- ParentResourceID 
    '550e8400-e29b-41d4-a716-446655440005', -- ResourceTypeID (Account resource type ID)
    'Account_Organization',                   -- Name
    1,                                     -- RowStatus
    'admin_user',                          -- CreatedBy
    'admin_user',                          -- UpdatedBy
    '2024-01-01 10:00:00',                 -- CreatedAt
    '2024-01-01 10:00:00'                  -- UpdatedAt
);

-- Client
INSERT INTO tnt_resource (
    resource_id,
    parent_resource_id,
    resource_type_id,
    name,
    row_status,
    created_by,
    updated_by,
    created_at,
    updated_at
) VALUES
(
    '550e8400-1111-41d4-a716-446655442222', -- ResourceID
    '11111111-1111-1111-1111-111111111111', -- ParentResourceID 
    '550e8400-e29b-41d4-a716-446655440006', -- ResourceTypeID (Client resource type ID)
    'Client_Organization',                   -- Name
    1,                                     -- RowStatus
    'admin_user',                          -- CreatedBy
    'admin_user',                          -- UpdatedBy
    '2024-01-01 10:00:00',                 -- CreatedAt
    '2024-01-01 10:00:00'                  -- UpdatedAt
);