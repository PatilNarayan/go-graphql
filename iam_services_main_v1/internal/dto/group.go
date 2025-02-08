package dto

import "time"

type GroupEntity struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	TenantID  int       `json:"tenant_id"`
	RowStatus int       `json:"row_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupInput struct {
	Name     string `json:"name" binding:"required" validate:"min=3,max=100"`
	TenantID int    `json:"tenantId" binding:"required"`
}
