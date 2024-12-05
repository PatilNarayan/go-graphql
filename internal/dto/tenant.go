package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID            string         `gorm:"type:char(36);primaryKey" json:"id"` // Remove default
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	ParentOrgID   string         `gorm:"type:varchar(255)" json:"parent_org_id"`
	ContactInfoID string         `gorm:"type:varchar(255)" json:"contact_info_id"`
	RowStatus     int            `gorm:"default:1" json:"row_status"`
	Description   string         `gorm:"type:text" json:"description"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID before saving
func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
