-- Insert dummy data into mst_roles
INSERT INTO mst_roles (
    role_id, name, version, row_status, created_by, updated_by, created_at, updated_at
) VALUES
    ('bfe2d3a4-fc1a-4b7d-b4a3-4e8f44fba0f1', 'Admin', 'v1', 1, 'system', 'system', NOW(), NOW()),
    ('a8e7f2d6-3d2f-49f6-bcf1-1d0b7a70fcd2', 'Editor', 'v1', 1, 'system', 'system', NOW(), NOW()),
    ('d0f1e7b9-cf2d-4d6a-9e1f-7b2a4a6d9c3e', 'Viewer', 'v1', 1, 'system', 'system', NOW(), NOW());

-- Insert dummy data into mst_permissions
INSERT INTO mst_permissions (
    permission_id, service_id, name, action, row_status, created_by, updated_by, created_at, updated_at
) VALUES
    ('f1e3d9a7-7d6c-4f2a-9b1e-2a3d7c6b9f4e', 'a2d6e3b7-cf1d-4f9a-b7d2-1f3a7e6c9b11', 'View Dashboard', 'view_dashboard', 1, 'system', 'system', NOW(), NOW()),
    ('a2d6e3b7-cf1d-4f9a-b7d2-1f3a7e6c9b4f', 'a2d6e3b7-cf1d-4f9a-b7d2-1f3a7e6c9b12', 'Edit Dashboard', 'edit_dashboard', 1, 'system', 'system', NOW(), NOW()),
    ('b7d9e1f2-a6d3-4f7a-b2d6-9f1a3e7c4b2f', 'a2d6e3b7-cf1d-4f9a-b7d2-1f3a7e6c9b13', 'Manage Users', 'manage_users', 1, 'system', 'system', NOW(), NOW()),
    ('d6f9e3a2-7b1f-4c6a-b9d7-2f3a6e1c9b7f', 'a2d6e3b7-cf1d-4f9a-b7d2-1f3a7e6c9b14', 'Delete Records', 'delete_records', 1, 'system', 'system', NOW(), NOW());

-- Insert dummy data into mst_role_permissions
INSERT INTO mst_role_permissions (
    role_permission_id, role_id, permission_id, row_status, created_by, updated_by, created_at
) VALUES
    ('c9e1d7f3-b6a4-4f9a-b2d6-7f3a9e1c4b7f', 'bfe2d3a4-fc1a-4b7d-b4a3-4e8f44fba0f1', 'f1e3d9a7-7d6c-4f2a-9b1e-2a3d7c6b9f4e', 1, 'system', 'system', NOW()),
    ('e3b9f1a7-d6c4-4a7b-b2f9-1a7d6e3c4b9f', 'bfe2d3a4-fc1a-4b7d-b4a3-4e8f44fba0f1', 'a2d6e3b7-cf1d-4f9a-b7d2-1f3a7e6c9b4f', 1, 'system', 'system', NOW()),
    ('f7d2b9e1-a6f3-4d9a-b1d6-7e9a3f1c4b7f', 'bfe2d3a4-fc1a-4b7d-b4a3-4e8f44fba0f1', 'b7d9e1f2-a6d3-4f7a-b2d6-9f1a3e7c4b2f', 1, 'system', 'system', NOW()),
    ('a6b9f3d7-2e1f-4c7a-b1d9-3f7a6e9c4b1f', 'a8e7f2d6-3d2f-49f6-bcf1-1d0b7a70fcd2', 'f1e3d9a7-7d6c-4f2a-9b1e-2a3d7c6b9f4e', 1, 'system', 'system', NOW()),
    ('b9e7a6f3-d1f2-4c9a-b7f6-1d3e9a7c4b2f', 'a8e7f2d6-3d2f-49f6-bcf1-1d0b7a70fcd2', 'a2d6e3b7-cf1d-4f9a-b7d2-1f3a7e6c9b4f', 1, 'system', 'system', NOW()),
    ('c7d9f2a6-b1e3-4f7a-b2d6-9e3a7f1c4b2f', 'd0f1e7b9-cf2d-4d6a-9e1f-7b2a4a6d9c3e', 'f1e3d9a7-7d6c-4f2a-9b1e-2a3d7c6b9f4e', 1, 'system', 'system', NOW());
