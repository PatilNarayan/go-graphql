ALTER TABLE tenants DROP COLUMN metadata;

ALTER TABLE tenants ALTER COLUMN metadata SET DEFAULT '{}';

SELECT * FROM root.tenants;