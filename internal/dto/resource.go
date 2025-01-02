package dto

import (
	"time"

	"github.com/google/uuid"
)

type TenantResource struct {
	ResourceID       uuid.UUID  `gorm:"type:char(36);primaryKey;column:resource_id" json:"resource_id"`
	ParentResourceID *uuid.UUID `gorm:"type:char(36);column:parent_resource_id" json:"parent_resource_id"`
	ResourceTypeID   uuid.UUID  `gorm:"type:char(36);not null;column:resource_type_id" json:"resource_type_id"` // foreign key to resource_type
	Name             string     `gorm:"size:45;not null;column:name" json:"name"`
	RowStatus        int        `gorm:"default:1" json:"row_status"`
	CreatedBy        string     `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy        string     `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *TenantResource) TableName() string {
	return "tnt_resource"
}

type Mst_ResourceType struct {
	ResourceTypeID uuid.UUID `gorm:"type:char(36);primaryKey;column:resource_type_id" json:"resource_type_id"`
	ServiceID      uuid.UUID `gorm:"type:char(36);not null;column:service_id" json:"service_id"`
	Name           string    `gorm:"size:45;not null;column:name" json:"name"`
	RowStatus      int       `gorm:"default:1" json:"row_status"`
	CreatedBy      string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy      string    `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *Mst_ResourceType) TableName() string {
	return "mst_resource_type"
}
