// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type Address struct {
	ID      string  `json:"id"`
	Street  *string `json:"street,omitempty"`
	City    *string `json:"city,omitempty"`
	State   *string `json:"state,omitempty"`
	ZipCode *string `json:"zipCode,omitempty"`
	Country *string `json:"country,omitempty"`
}

type AddressInput struct {
	Street  *string `json:"street,omitempty"`
	City    *string `json:"city,omitempty"`
	State   *string `json:"state,omitempty"`
	ZipCode *string `json:"zipCode,omitempty"`
	Country *string `json:"country,omitempty"`
}

type ContactInfo struct {
	ID          string   `json:"id"`
	Email       *string  `json:"email,omitempty"`
	PhoneNumber *string  `json:"phoneNumber,omitempty"`
	Address     *Address `json:"address,omitempty"`
}

type GroupInput struct {
	Name     string `json:"name"`
	TenantID int    `json:"tenantId"`
}

type TenantInput struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	ParentOrgID   string  `json:"parentOrgId"`
	ContactInfoID string  `json:"contactInfoId"`
}
