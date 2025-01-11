package dto

import (
	"time"

	"github.com/google/uuid"
)

type TNTRole struct {
	ResourceID       uuid.UUID  `json:"resourceId" gorm:"type:char(36);primaryKey;column:resource_id" db:"resource_id"`
	RoleType         string     `json:"roleType" gorm:"column:role_type;type:varchar(255);not null" db:"role_type"`
	ParentResourceID *uuid.UUID `json:"parentResourceId" gorm:"type:char(36);column:parent_resource_id" db:"parent_resource_id"`
	Name             string     `json:"name" gorm:"column:name;size:255" db:"name"`
	Version          string     `json:"version" gorm:"column:version;size:100" db:"version"`
	Description      string     `json:"description" gorm:"column:description;type:text" db:"description"`
	RowStatus        int        `json:"rowStatus" gorm:"column:row_status" db:"row_status"`
	CreatedBy        string     `json:"createdBy" gorm:"column:created_by;size:36" db:"created_by"`
	UpdatedBy        string     `json:"updatedBy" gorm:"column:updated_by;size:36" db:"updated_by"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
}

// TableName overrides the default table name
func (TNTRole) TableName() string {
	return "tnt_roles"
}

type TNTPermission struct {
	PermissionID uuid.UUID `json:"permissionId" gorm:"type:char(36);primaryKey;column:permission_id" db:"permission_id"`
	RoleID       uuid.UUID `json:"roleId" gorm:"type:char(36);column:role_id" db:"role_id"`
	ServiceID    string    `json:"serviceId" gorm:"column:service_id;size:36" db:"service_id"`
	Name         string    `json:"name" gorm:"column:name;size:255" db:"name"`
	Action       string    `json:"action" gorm:"column:action;size:100" db:"action"`
	RowStatus    int       `json:"rowStatus" gorm:"column:row_status" db:"row_status"`
	CreatedBy    string    `json:"createdBy" gorm:"column:created_by;size:36" db:"created_by"`
	UpdatedBy    string    `json:"updatedBy" gorm:"column:updated_by;size:36" db:"updated_by"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
}

// TableName overrides the default table name
func (TNTPermission) TableName() string {
	return "tnt_permissions"
}
