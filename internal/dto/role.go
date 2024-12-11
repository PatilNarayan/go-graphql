package dto

import (
	"time"
)

type Role struct {
	RoleID      string    `gorm:"type:char(36);primaryKey;column:role_id" json:"role_id"` // Use char(36) for UUID
	Name        string    `gorm:"size:45;not null;column:name" json:"name"`
	Description string    `gorm:"size:255;column:description" json:"description,omitempty"`
	Version     string    `gorm:"size:45;not null;column:version" json:"version"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	RoleType    string    `gorm:"type:enum('DEFAULT', 'CUSTOM');not null;column:role_type" json:"role_type"`
}

type Permission struct {
	PermissionID string    `gorm:"type:char(36);primaryKey;column:permission_id" json:"permission_id"` // Use char(36) for UUID
	Name         string    `gorm:"size:45;not null;column:name" json:"name"`
	Description  string    `gorm:"size:255;column:description" json:"description,omitempty"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	CreatedBy    string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedDate  time.Time `gorm:"column:updated_date;autoUpdateTime" json:"updated_date"`
	UpdatedBy    string    `gorm:"size:45;column:updated_by" json:"updated_by"`
}

type RolePermission struct {
	RolePermissionID string    `gorm:"type:char(36);primaryKey;column:role_permission_id" json:"role_permission_id"` // Use char(36) for UUID
	RoleID           string    `gorm:"type:char(36);not null;column:role_id" json:"role_id"`                         // Use char(36) for UUID
	PermissionID     string    `gorm:"type:char(36);not null;column:permission_id" json:"permission_id"`             // Use char(36) for UUID
	CreatedDate      time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	CreatedBy        string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedDate      time.Time `gorm:"column:updated_date;autoUpdateTime" json:"updated_date"`
	UpdatedBy        string    `gorm:"size:45;column:updated_by" json:"updated_by"`
}

type RoleScope struct {
	RoleScopeID    string    `gorm:"type:char(36);primaryKey;column:role_scope_id" json:"role_scope_id"`     // Use char(36) for UUID
	RoleID         string    `gorm:"type:char(36);not null;column:role_id" json:"role_id"`                   // Use char(36) for UUID
	ResourceTypeID string    `gorm:"type:char(36);not null;column:resource_type_id" json:"resource_type_id"` // Use char(36) for UUID
	CreatedDate    time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	CreatedBy      string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedDate    time.Time `gorm:"column:updated_date;autoUpdateTime" json:"updated_date"`
	UpdatedBy      string    `gorm:"size:45;column:updated_by" json:"updated_by"`
}
