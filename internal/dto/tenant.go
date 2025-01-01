package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tenant struct aligned with schema
type TenantMetadata struct {
	ID         uuid.UUID       `gorm:"type:char(36);primaryKey;column:id" json:"id"`
	ResourceID string          `gorm:"type:char(36);not null" json:"resource_id"`
	Metadata   json.RawMessage `gorm:"type:json;" json:"metadata"`
	RowStatus  int             `gorm:"default:1" json:"row_status"`
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
