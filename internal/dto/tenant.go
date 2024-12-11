package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type ContactInfo struct {
// 	ID          string `gorm:"type:char(36);primaryKey" json:"id"`
// 	Email       string `gorm:"type:varchar(255)" json:"email"`
// 	PhoneNumber string `gorm:"type:varchar(15)" json:"phone_number"`
// 	Address     string `gorm:"type:text" json:"address"` // Address could reference a detailed Address model if needed
// }

// Tenant struct aligned with schema
type Tenant struct {
	ID            string `gorm:"type:char(36);primaryKey" json:"id"` // UUID
	Name          string `gorm:"type:varchar(255);not null" json:"name"`
	ParentOrgID   string `gorm:"type:char(36)" json:"parent_org_id"`
	ContactInfoID string `gorm:"type:char(36)" json:"contact_info_id"`
	// ContactInfo   ContactInfo    `gorm:"foreignKey:ContactInfoID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"contact_info"`
	RowStatus   int            `gorm:"default:1" json:"row_status"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
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
