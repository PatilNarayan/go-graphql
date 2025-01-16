-- Create indexes for better performance
CREATE INDEX idx_mst_resource_types_service_id ON mst_resource_types(service_id);
CREATE INDEX idx_mst_resource_types_name ON mst_resource_types(name);

-- Insert dummy data for Root resource into tnt_resource table
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
    '11111111-1111-1111-1111-111111111111', -- ResourceID
    NULL,                                  -- ParentResourceID (Root has no parent)
    '550e8400-e29b-41d4-a716-446655440004', -- ResourceTypeID (Root resource type ID)
    'Root Organization',                   -- Name
    1,                                     -- RowStatus
    'admin_user',                          -- CreatedBy
    'admin_user',                          -- UpdatedBy
    '2024-01-01 10:00:00',                 -- CreatedAt
    '2024-01-01 10:00:00'                  -- UpdatedAt
);

-- Insert dummy data for Root user into tnt_resource_metadata table
INSERT INTO tnt_resource_metadata (
    id,
    resource_id,
    metadata,
    row_status,
    created_by,
    updated_by,
    created_at,
    updated_at,
    deleted_at
) VALUES 
(
    '22222222-2222-2222-2222-222222222222', -- ID (UUID for the metadata entry)
    '11111111-1111-1111-1111-111111111111', -- ResourceID (references the root resource ID in tnt_resource)
    '{"description": "Root organization metadata", "contactInfo": {"email": "root@organization.com", "phone": "1234567890"}}', -- Metadata (JSON format)
    1,                                      -- RowStatus (active)
    'admin_user',                           -- CreatedBy
    'admin_user',                           -- UpdatedBy
    '2024-01-01 10:00:00',                  -- CreatedAt
    '2024-01-01 10:00:00',                  -- UpdatedAt
    NULL                                    -- DeletedAt (NULL for active records)
);
