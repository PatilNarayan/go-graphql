package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleTypeEnum string

const (
	RoleTypeEnumDefault RoleTypeEnum = "DEFAULT"
	RoleTypeEnumCustom  RoleTypeEnum = "CUSTOM"
)

type TNTRole struct {
	ResourceID     uuid.UUID      `json:"resourceId" gorm:"type:char(36);primaryKey;column:resource_id" db:"resource_id"`
	RoleType       RoleTypeEnum   `json:"roleType" gorm:"column:role_type" db:"role_type"`
	Name           string         `json:"name" gorm:"column:name;size:255" db:"name"`
	Version        string         `json:"version" gorm:"column:version;size:100" db:"version"`
	ResourceTypeID uuid.UUID      `json:"resourceTypeId" gorm:"type:char(36);not null;column:resource_type_id" db:"resource_type_id"`
	Description    string         `json:"description" gorm:"column:description;type:text" db:"description"`
	RowStatus      int            `json:"rowStatus" gorm:"column:row_status" db:"row_status"`
	CreatedBy      string         `json:"createdBy" gorm:"column:created_by;size:36" db:"created_by"`
	UpdatedBy      string         `json:"updatedBy" gorm:"column:updated_by;size:36" db:"updated_by"`
	CreatedAt      time.Time      `json:"createdAt" gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName overrides the default table name
func (TNTRole) TableName() string {
	return "tnt_roles"
}

type TNTPermission struct {
	PermissionID uuid.UUID      `json:"permissionId" gorm:"type:char(36);primaryKey;column:permission_id" db:"permission_id"`
	ServiceID    string         `json:"serviceId" gorm:"column:service_id;size:36" db:"service_id"`
	Name         string         `json:"name" gorm:"column:name;size:255" db:"name"`
	Action       string         `json:"action" gorm:"column:action;size:100" db:"action"`
	RowStatus    int            `json:"rowStatus" gorm:"column:row_status" db:"row_status"`
	CreatedBy    string         `json:"createdBy" gorm:"column:created_by;size:36" db:"created_by"`
	UpdatedBy    string         `json:"updatedBy" gorm:"column:updated_by;size:36" db:"updated_by"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName overrides the default table name
func (TNTPermission) TableName() string {
	return "tnt_permissions"
}

// RolePermission represents the many-to-many relationship between roles and permissions
type TNTRolePermission struct {
	ID           uuid.UUID      `gorm:"column:role_permission_id;type:varchar(36);primary_key" json:"id"`
	RoleID       uuid.UUID      `gorm:"column:role_id;type:varchar(36);not null" json:"roleId"`
	PermissionID uuid.UUID      `gorm:"column:permission_id;type:varchar(36);not null" json:"permissionId"`
	RowStatus    int            `gorm:"column:row_status;type:tinyint(1);default:1" json:"rowStatus"`
	CreatedBy    string         `gorm:"column:created_by;type:varchar(36)" json:"createdBy"`
	UpdatedBy    string         `gorm:"column:updated_by;type:varchar(36)" json:"updatedBy"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the RolePermission model
func (TNTRolePermission) TableName() string {
	return "tnt_role_permissions"
}

type MstRole struct {
	RoleID         uuid.UUID      `json:"roleId" gorm:"type:char(36);primaryKey;column:role_id" db:"role_id"`
	Name           string         `json:"name" gorm:"column:name;size:255" db:"name"`
	Version        string         `json:"version" gorm:"column:version;size:100" db:"version"`
	ResourceTypeID uuid.UUID      `json:"resourceTypeId" gorm:"type:char(36);not null;column:resource_type_id" db:"resource_type_id"`
	Description    string         `json:"description" gorm:"column:description;type:text" db:"description"`
	RowStatus      int            `json:"rowStatus" gorm:"column:row_status" db:"row_status"`
	CreatedBy      string         `json:"createdBy" gorm:"column:created_by;size:36" db:"created_by"`
	UpdatedBy      string         `json:"updatedBy" gorm:"column:updated_by;size:36" db:"updated_by"`
	CreatedAt      time.Time      `json:"createdAt" gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

func (MstRole) TableName() string {
	return "mst_roles"
}

type MstPermission struct {
	PermissionID uuid.UUID      `json:"permissionId" gorm:"type:char(36);primaryKey;column:permission_id" db:"permission_id"`
	ServiceID    string         `json:"serviceId" gorm:"column:service_id;size:36" db:"service_id"`
	Name         string         `json:"name" gorm:"column:name;size:255" db:"name"`
	Action       string         `json:"action" gorm:"column:action;size:100" db:"action"`
	RowStatus    int            `json:"rowStatus" gorm:"column:row_status" db:"row_status"`
	CreatedBy    string         `json:"createdBy" gorm:"column:created_by;size:36" db:"created_by"`
	UpdatedBy    string         `json:"updatedBy" gorm:"column:updated_by;size:36" db:"updated_by"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

func (MstPermission) TableName() string {
	return "mst_permissions"
}

type MstRolePermission struct {
	ID           uuid.UUID      `gorm:"column:role_permission_id;type:varchar(36);primary_key" json:"id"`
	RoleID       uuid.UUID      `gorm:"column:role_id;type:varchar(36);not null" json:"roleId"`
	PermissionID uuid.UUID      `gorm:"column:permission_id;type:varchar(36);not null" json:"permissionId"`
	RowStatus    int            `gorm:"column:row_status;type:tinyint(1);default:1" json:"rowStatus"`
	CreatedBy    string         `gorm:"column:created_by;type:varchar(36)" json:"createdBy"`
	UpdatedBy    string         `gorm:"column:updated_by;type:varchar(36)" json:"updatedBy"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (MstRolePermission) TableName() string {
	return "mst_role_permissions"
}
