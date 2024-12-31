// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"github.com/google/uuid"
)

type Organization interface {
	IsOrganization()
	GetID() uuid.UUID
	GetName() string
	GetDescription() *string
	GetCreatedAt() string
	GetUpdatedAt() *string
	GetCreatedBy() *string
	GetUpdatedBy() *string
}

type Resource interface {
	IsResource()
	GetID() uuid.UUID
	GetName() string
	GetCreatedAt() string
	GetUpdatedAt() *string
	GetCreatedBy() *string
	GetUpdatedBy() *string
}

type Address struct {
	ID      uuid.UUID `json:"id"`
	Street  *string   `json:"street,omitempty"`
	City    *string   `json:"city,omitempty"`
	State   *string   `json:"state,omitempty"`
	ZipCode *string   `json:"zipCode,omitempty"`
	Country *string   `json:"country,omitempty"`
}

type ClientOrganizationUnit struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Tenant      *Tenant      `json:"tenant"`
	ParentOrg   Organization `json:"parentOrg"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   *string      `json:"updated_at,omitempty"`
	CreatedBy   *string      `json:"created_by,omitempty"`
	UpdatedBy   *string      `json:"updated_by,omitempty"`
}

func (ClientOrganizationUnit) IsResource()                {}
func (this ClientOrganizationUnit) GetID() uuid.UUID      { return this.ID }
func (this ClientOrganizationUnit) GetName() string       { return this.Name }
func (this ClientOrganizationUnit) GetCreatedAt() string  { return this.CreatedAt }
func (this ClientOrganizationUnit) GetUpdatedAt() *string { return this.UpdatedAt }
func (this ClientOrganizationUnit) GetCreatedBy() *string { return this.CreatedBy }
func (this ClientOrganizationUnit) GetUpdatedBy() *string { return this.UpdatedBy }

func (ClientOrganizationUnit) IsOrganization() {}

func (this ClientOrganizationUnit) GetDescription() *string { return this.Description }

type ContactInfo struct {
	Email       *string  `json:"email,omitempty"`
	PhoneNumber *string  `json:"phoneNumber,omitempty"`
	Address     *Address `json:"address,omitempty"`
}

type ContactInfoInput struct {
	Email       *string             `json:"email,omitempty"`
	PhoneNumber *string             `json:"phoneNumber,omitempty"`
	Address     *CreateAddressInput `json:"address,omitempty"`
}

type CreateAddressInput struct {
	Street  *string `json:"street,omitempty"`
	City    *string `json:"city,omitempty"`
	State   *string `json:"state,omitempty"`
	ZipCode *string `json:"zipCode,omitempty"`
	Country *string `json:"country,omitempty"`
}

type CreateClientOrganizationUnitInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	TenantID    string  `json:"tenantId"`
	ParentOrgID string  `json:"parentOrgId"`
}

type CreateRootInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type CreateTenantInput struct {
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	ParentOrgID string            `json:"parentOrgId"`
	ContactInfo *ContactInfoInput `json:"contactInfo,omitempty"`
}

type Root struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   *string   `json:"updated_at,omitempty"`
	CreatedBy   *string   `json:"created_by,omitempty"`
	UpdatedBy   *string   `json:"updated_by,omitempty"`
}

func (Root) IsResource()                {}
func (this Root) GetID() uuid.UUID      { return this.ID }
func (this Root) GetName() string       { return this.Name }
func (this Root) GetCreatedAt() string  { return this.CreatedAt }
func (this Root) GetUpdatedAt() *string { return this.UpdatedAt }
func (this Root) GetCreatedBy() *string { return this.CreatedBy }
func (this Root) GetUpdatedBy() *string { return this.UpdatedBy }

func (Root) IsOrganization() {}

func (this Root) GetDescription() *string { return this.Description }

type Tenant struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	ParentOrg   Organization `json:"parentOrg"`
	ContactInfo *ContactInfo `json:"contactInfo,omitempty"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   *string      `json:"updated_at,omitempty"`
	CreatedBy   *string      `json:"created_by,omitempty"`
	UpdatedBy   *string      `json:"updated_by,omitempty"`
}

func (Tenant) IsResource()                {}
func (this Tenant) GetID() uuid.UUID      { return this.ID }
func (this Tenant) GetName() string       { return this.Name }
func (this Tenant) GetCreatedAt() string  { return this.CreatedAt }
func (this Tenant) GetUpdatedAt() *string { return this.UpdatedAt }
func (this Tenant) GetCreatedBy() *string { return this.CreatedBy }
func (this Tenant) GetUpdatedBy() *string { return this.UpdatedBy }

func (Tenant) IsOrganization() {}

func (this Tenant) GetDescription() *string { return this.Description }

type UpdateAddressInput struct {
	Street  *string `json:"street,omitempty"`
	City    *string `json:"city,omitempty"`
	State   *string `json:"state,omitempty"`
	ZipCode *string `json:"zipCode,omitempty"`
	Country *string `json:"country,omitempty"`
}

type UpdateClientOrganizationUnitInput struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	TenantID    *string   `json:"tenantId,omitempty"`
	ParentOrgID *string   `json:"parentOrgId,omitempty"`
}

type UpdateRootInput struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
}

type UpdateTenantInput struct {
	ID            uuid.UUID `json:"id"`
	Name          *string   `json:"name,omitempty"`
	Description   *string   `json:"description,omitempty"`
	ParentOrgID   *string   `json:"parentOrgId,omitempty"`
	ContactInfoID string    `json:"contactInfoId"`
}
