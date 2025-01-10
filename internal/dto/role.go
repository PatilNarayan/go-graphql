package dto

import "time"

type TNTRole struct {
	ResourceID string    `json:"resourceId" gorm:"column:resource_id;primaryKey;size:36" db:"resource_id"`
	RoleType   string    `json:"roleType" gorm:"column:role_type;type:varchar(255);not null" db:"role_type"`
	Name       string    `json:"name" gorm:"column:name;size:255" db:"name"`
	Version    string    `json:"version" gorm:"column:version;size:100" db:"version"`
	RowStatus  bool      `json:"rowStatus" gorm:"column:row_status" db:"row_status"`
	CreatedBy  string    `json:"createdBy" gorm:"column:created_by;size:36" db:"created_by"`
	UpdatedBy  string    `json:"updatedBy" gorm:"column:updated_by;size:36" db:"updated_by"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
}

// TableName overrides the default table name
func (TNTRole) TableName() string {
	return "tnt_roles"
}

type TNTPermission struct {
	PermissionID string    `json:"permissionId" gorm:"column:permission_id;primaryKey;size:36" db:"permission_id"`
	ServiceID    string    `json:"serviceId" gorm:"column:service_id;size:36" db:"service_id"`
	Name         string    `json:"name" gorm:"column:name;size:255" db:"name"`
	Action       string    `json:"action" gorm:"column:action;size:100" db:"action"`
	RowStatus    bool      `json:"rowStatus" gorm:"column:row_status" db:"row_status"`
	CreatedBy    string    `json:"createdBy" gorm:"column:created_by;size:36" db:"created_by"`
	UpdatedBy    string    `json:"updatedBy" gorm:"column:updated_by;size:36" db:"updated_by"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
}

// TableName overrides the default table name
func (TNTPermission) TableName() string {
	return "tnt_permissions"
}
