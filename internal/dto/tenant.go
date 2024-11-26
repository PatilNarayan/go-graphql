package dto

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"` // Auto-generate UUID
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`                    // Name column with length restriction
	ParentOrgID    string         `gorm:"type:varchar(255)" json:"parent_org_id"`                    // Parent organization ID // Not null constraint
	ContactInfoID  string         `gorm:"type:varchar(255)" json:"contact_info_id"`                  // New contact info ID column
	RowStatus      int            `gorm:"default:1" json:"row_status"`                               // Default value for row status
	Description    string         `gorm:"type:text" json:"description"`                              // New description column
	RemoteTenantID string         `gorm:"type:varchar(255)" json:"remote_tenant_id"`                 // Remote tenant ID with column type
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`                          // Auto-managed by GORM
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                          // Auto-managed by GORM
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`                                            // For soft deletes (optional)
}
