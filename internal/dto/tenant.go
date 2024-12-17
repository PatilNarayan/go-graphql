package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tenant struct aligned with schema
type Tenant struct {
	ID             string         `gorm:"type:char(36);primaryKey" json:"id"`
	ResourceID     string         `gorm:"type:char(36);not null" json:"resource_id"`
	ParentTenantID string         `gorm:"type:char(36)" json:"parent_tenant_id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	ParentOrgID    string         `gorm:"type:char(36)" json:"parent_org_id"`
	ContactInfoID  string         `gorm:"type:char(36)" json:"contact_info_id"`
	RowStatus      int            `gorm:"default:1" json:"row_status"`
	Description    string         `gorm:"type:text" json:"description"`
	Metadata       string         `gorm:"type:text" json:"metadata"`
	CreatedBy      string         `gorm:"size:45" json:"created_by"`
	UpdatedBy      string         `gorm:"size:45" json:"updated_by"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// BeforeCreate hook to generate UUID before saving
func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}

// // BeforeCreate hook for ContactInfo
// func (c *ContactInfo) BeforeCreate(tx *gorm.DB) (err error) {
// 	c.ID = uuid.New().String()
// 	return
// }
