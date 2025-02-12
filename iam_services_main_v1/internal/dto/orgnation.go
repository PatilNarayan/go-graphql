package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type Organization struct {
	OrganizationID uuid.UUID `gorm:"type:char(36);primaryKey;column:organization_id" json:"organization_id"`
	Name           string    `gorm:"size:36;not null;column:name" json:"name"`
	Description    string    `gorm:"type:text;column:description" json:"description"`
	ParentOrgId    uuid.UUID `gorm:"type:char(36);column:parent_org_id" json:"parent_org_id"`
	RowStatus      int       `gorm:"default:1" json:"row_status"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
