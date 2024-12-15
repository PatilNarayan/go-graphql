package dto

import (
	"time"
)

type Role struct {
	RoleID      string    `gorm:"type:char(36);primaryKey;column:role_id" json:"role_id"` // Use char(36) for UUID
	ResourceID  string    `gorm:"type:char(36);not null;column:resource_id" json:"resource_id"`
	Name        string    `gorm:"size:45;not null;column:name" json:"name"`
	RoleType    string    `gorm:"type:text;not null;column:role_type;check:role_type IN ('DEFAULT', 'CUSTOM')" json:"role_type"`
	RowStatus   int       `gorm:"default:1" json:"row_status"`
	Description string    `gorm:"type:text;column:description" json:"description"`
	Version     string    `gorm:"column:version" json:"version"`
	CreatedBy   string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy   string    `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (r *Role) TableName() string {
	return "role"
}

type Permission struct {
	PermissionID string    `gorm:"type:char(36);primaryKey;column:permission_id" json:"permission_id"` // Use char(36) for UUID
	Name         string    `gorm:"size:45;not null;column:name" json:"name"`
	RowStatus    int       `gorm:"default:1" json:"row_status"`
	ServiceID    string    `gorm:"size:45;column:service_id" json:"service_id"`
	Action       string    `gorm:"size:45;column:action" json:"action"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"created_date"`
	CreatedBy    string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedDate  time.Time `gorm:"column:updated_date;autoUpdateTime" json:"updated_date"`
	UpdatedBy    string    `gorm:"size:45;column:updated_by" json:"updated_by"`
}

type RoleAssignment struct {
	RoleAssignmentID string    `gorm:"type:char(36);primaryKey;column:role_assignment_id" json:"role_assignment_id"` // Use char(36) for UUID
	PrincipalID      string    `gorm:"type:char(36);not null;column:principal_id" json:"principal_id"`
	RoleID           string    `gorm:"type:char(36);not null;column:role_id" json:"role_id"`
	PermissionID     string    `gorm:"type:char(36);not null;column:permission_id" json:"permission_id"`
	RowStatus        int       `gorm:"default:1" json:"row_status"`
	CreatedBy        string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy        string    `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Updated_at       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

type Principal struct {
	PrincipalID string    `gorm:"type:char(36);primaryKey;column:principal_id" json:"principal_id"` // Use char(36) for UUID
	TenantID    string    `gorm:"type:char(36);not null;column:tenant_id" json:"tenant_id"`
	Name        string    `gorm:"size:45;not null;column:name" json:"name"`
	Email       string    `gorm:"size:45;not null;column:email" json:"email"`
	RowStatus   int       `gorm:"default:1" json:"row_status"`
	CreatedBy   string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy   string    `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Updated_at  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
