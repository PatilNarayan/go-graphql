package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantResource struct {
	ResourceID       uuid.UUID      `gorm:"type:char(36);primaryKey;column:resource_id" json:"resource_id"`
	ParentResourceID *uuid.UUID     `gorm:"type:char(36);column:parent_resource_id" json:"parent_resource_id"`
	ResourceTypeID   uuid.UUID      `gorm:"type:char(36);not null;column:resource_type_id" json:"resource_type_id"` // foreign key to resource_type
	Name             string         `gorm:"size:45;not null;column:name" json:"name"`
	TenantID         *uuid.UUID     `gorm:"type:char(36);column:tenant_id" json:"tenant_id"`
	RowStatus        int            `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy        string         `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy        string         `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

func (t *TenantResource) TableName() string {
	return "tnt_resource"
}

type Mst_ResourceTypes struct {
	ResourceTypeID uuid.UUID `gorm:"type:char(36);primaryKey;column:resource_type_id" json:"resource_type_id"`
	ServiceID      uuid.UUID `gorm:"type:char(36);not null;column:service_id" json:"service_id"`
	Name           string    `gorm:"size:45;not null;column:name" json:"name"`
	RowStatus      int       `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy      string    `gorm:"size:45;column:created_by" json:"created_by"`
	UpdatedBy      string    `gorm:"size:45;column:updated_by" json:"updated_by"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

func (t *Mst_ResourceTypes) TableName() string {
	return "mst_resource_types"
}

// Tenant struct aligned with schema
type TenantMetadata struct {
	ID         uuid.UUID       `gorm:"type:char(36);primaryKey;column:id" json:"id"`
	ResourceID string          `gorm:"type:char(36);not null" json:"resource_id"`
	Metadata   json.RawMessage `gorm:"type:json;" json:"metadata"`
	RowStatus  int             `gorm:"default:1;column:row_status" json:"row_status"`
	CreatedBy  string          `gorm:"size:45" json:"created_by"`
	UpdatedBy  string          `gorm:"size:45" json:"updated_by"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"` // Soft delete
}

// BeforeCreate hook to generate UUID before saving
func (t *TenantMetadata) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

func (t *TenantMetadata) TableName() string {
	return "tnt_resource_metadata"
}
