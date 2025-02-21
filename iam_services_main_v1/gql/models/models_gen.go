// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
)

// Define a union for the possible 'data' types
type Data interface {
	IsData()
}

// Standard Error Interface for the error responses
type Error interface {
	IsError()
	// Error code representing the type of error.
	GetErrorCode() string
	// Details about the error.
	GetErrorDetails() *string
	// A message providing information about the operation to the user.
	GetMessage() string
	// A message providing additional context or information about the operation for the logging.
	GetSystemMessage() string
}

// Define a union for the possible operation results
type OperationResult interface {
	IsOperationResult()
}

// Represents an Organization entity
type Organization interface {
	IsOrganization()
	// Timestamp of creation
	GetCreatedAt() string
	// Identifier of the user who created the record
	GetCreatedBy() uuid.UUID
	// Description of the organization
	GetDescription() *string
	// Unique identifier of the organization
	GetID() uuid.UUID
	// Name of the organization
	GetName() string
	// Parent organization
	GetParentOrg() Organization
	// Timestamp of last update
	GetUpdatedAt() string
	// Identifier of the user who last updated the record
	GetUpdatedBy() uuid.UUID
}

// Represents a Principal entity
type Principal interface {
	IsPrincipal()
	// Email of the principal
	GetEmail() string
	// Unique identifier of the principal
	GetID() uuid.UUID
	// Name of the principal
	GetName() string
	// Tenant associated with the principal
	GetTenant() *Tenant
}

// Represents a Resource entity
type Resource interface {
	IsResource()
	// Timestamp of creation
	GetCreatedAt() string
	// Identifier of the user who created the record
	GetCreatedBy() uuid.UUID
	// Unique identifier of the resource
	GetID() uuid.UUID
	// Name of the resource
	GetName() string
	// Timestamp of last update
	GetUpdatedAt() string
	// Identifier of the user who last updated the record
	GetUpdatedBy() uuid.UUID
}

// Standard Response Interface for both success and error responses
type Response interface {
	IsResponse()
	// Indicates if the operation was successful.
	GetIsSuccess() bool
	// A message providing additional context or information about the operation.
	GetMessage() string
}

// Represents an Account entity
type Account struct {
	// Billing Info entity
	BillingInfo *BillingInfo `json:"billingInfo,omitempty"`
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Description of the account
	Description *string `json:"description,omitempty"`
	// Unique identifier of the account
	ID uuid.UUID `json:"id"`
	// Name of the account
	Name string `json:"name"`
	// Parent organization
	ParentOrg Organization `json:"parentOrg,omitempty"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

func (Account) IsData() {}

func (Account) IsOrganization() {}

// Timestamp of creation
func (this Account) GetCreatedAt() string { return this.CreatedAt }

// Identifier of the user who created the record
func (this Account) GetCreatedBy() uuid.UUID { return this.CreatedBy }

// Description of the organization
func (this Account) GetDescription() *string { return this.Description }

// Unique identifier of the organization
func (this Account) GetID() uuid.UUID { return this.ID }

// Name of the organization
func (this Account) GetName() string { return this.Name }

// Parent organization
func (this Account) GetParentOrg() Organization { return this.ParentOrg }

// Timestamp of last update
func (this Account) GetUpdatedAt() string { return this.UpdatedAt }

// Identifier of the user who last updated the record
func (this Account) GetUpdatedBy() uuid.UUID { return this.UpdatedBy }

func (Account) IsResource() {}

// Timestamp of creation

// Identifier of the user who created the record

// Unique identifier of the resource

// Name of the resource

// Timestamp of last update

// Identifier of the user who last updated the record

// Represents an address
type Address struct {
	// City of the address
	City *string `json:"city,omitempty"`
	// Country of the address
	Country *string `json:"country,omitempty"`
	// State of the address
	State *string `json:"state,omitempty"`
	// Street of the address
	Street *string `json:"street,omitempty"`
	// Zip code of the address
	ZipCode *string `json:"zipCode,omitempty"`
}

// Defines input fields for creating an address
type AddressInput struct {
	// City of the address
	City *string `json:"city,omitempty"`
	// Country of the address
	Country *string `json:"country,omitempty"`
	// State of the address
	State *string `json:"state,omitempty"`
	// Street of the address
	Street *string `json:"street,omitempty"`
	// Zip code of the address
	ZipCode *string `json:"zipCode,omitempty"`
}

// Represents a billing address entity associated to account
type BillingAddress struct {
	// Name of the city associated to billing address
	City string `json:"city"`
	// Name of the country associated to billing address
	Country string `json:"country"`
	// Name of the state associated to billing address
	State string `json:"state"`
	// Name of the street associated to billing address
	Street string `json:"street"`
	// Name of the zipcode associated to billing address
	Zipcode string `json:"zipcode"`
}

// Represents a billing info entity associated to account
type BillingInfo struct {
	// Billing Address associated to account
	BillingAddress *BillingAddress `json:"billingAddress"`
	// Credit card number associated to account
	CreditCardNumber string `json:"creditCardNumber"`
	// Credit card type associated to account
	CreditCardType string `json:"creditCardType"`
	// CVV associated to account
	Cvv string `json:"cvv"`
	// Expiration date associated to account
	ExpirationDate string `json:"expirationDate"`
}

// Represents a Binding entity
type Binding struct {
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Unique identifier of the binding
	ID uuid.UUID `json:"id"`
	// Name of the binding
	Name string `json:"name"`
	// Principal associated with the binding
	Principal Principal `json:"principal"`
	// Role associated with the binding
	Role *Role `json:"role"`
	// Scope reference associated with the binding
	ScopeRef Resource `json:"scopeRef"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
	// Version of the binding
	Version string `json:"version"`
}

func (Binding) IsData() {}

// Represents a Client Organization Unit entity
type ClientOrganizationUnit struct {
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Description of the client organization unit
	Description *string `json:"description,omitempty"`
	// Unique identifier of the client organization unit
	ID uuid.UUID `json:"id"`
	// Name of the client organization unit
	Name string `json:"name"`
	// Parent organization
	ParentOrg Organization `json:"parentOrg"`
	// Tenant associated with the client organization unit
	Tenant *Tenant `json:"tenant"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

func (ClientOrganizationUnit) IsData() {}

func (ClientOrganizationUnit) IsOrganization() {}

// Timestamp of creation
func (this ClientOrganizationUnit) GetCreatedAt() string { return this.CreatedAt }

// Identifier of the user who created the record
func (this ClientOrganizationUnit) GetCreatedBy() uuid.UUID { return this.CreatedBy }

// Description of the organization
func (this ClientOrganizationUnit) GetDescription() *string { return this.Description }

// Unique identifier of the organization
func (this ClientOrganizationUnit) GetID() uuid.UUID { return this.ID }

// Name of the organization
func (this ClientOrganizationUnit) GetName() string { return this.Name }

// Parent organization
func (this ClientOrganizationUnit) GetParentOrg() Organization { return this.ParentOrg }

// Timestamp of last update
func (this ClientOrganizationUnit) GetUpdatedAt() string { return this.UpdatedAt }

// Identifier of the user who last updated the record
func (this ClientOrganizationUnit) GetUpdatedBy() uuid.UUID { return this.UpdatedBy }

func (ClientOrganizationUnit) IsResource() {}

// Timestamp of creation

// Identifier of the user who created the record

// Unique identifier of the resource

// Name of the resource

// Timestamp of last update

// Identifier of the user who last updated the record

// Represents contact information
type ContactInfo struct {
	// Address of the contact
	Address *Address `json:"address,omitempty"`
	// Email of the contact
	Email *string `json:"email,omitempty"`
	// Phone number of the contact
	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

// Defines input fields for contact information
type ContactInfoInput struct {
	// Address of the contact
	Address *AddressInput `json:"address,omitempty"`
	// Email of the contact
	Email *string `json:"email,omitempty"`
	// Phone number of the contact
	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

// Defines input fields for creating an account
type CreateAccountInput struct {
	// Scope of billing info
	BillingInfo *CreateBillingInfoInput `json:"billingInfo,omitempty"`
	// Description of the account
	Description *string `json:"description,omitempty"`
	// Unique identifier of the account
	ID uuid.UUID `json:"id"`
	// Name of the account
	Name string `json:"name"`
	// Associated parent organization
	ParentID uuid.UUID `json:"parentId"`
	// Associated tenant
	TenantID uuid.UUID `json:"tenantId"`
}

// Defines input fields for creating a billing address for an account
type CreateBillingAddressInput struct {
	// Name of the city associated to billing address
	City string `json:"city"`
	// Name of the country associated to billing address
	Country string `json:"country"`
	// Name of the state associated to billing address
	State string `json:"state"`
	// Name of the street associated to billing address
	Street string `json:"street"`
	// Name of the zipcode associated to billing address
	Zipcode string `json:"zipcode"`
}

// Defines input fields for creating billing info for an account
type CreateBillingInfoInput struct {
	// Billing Address associated to account
	BillingAddress *CreateBillingAddressInput `json:"billingAddress"`
	// Credit card number associated to account
	CreditCardNumber string `json:"creditCardNumber"`
	// Credit card type associated to account
	CreditCardType string `json:"creditCardType"`
	// CVV associated to account
	Cvv string `json:"cvv"`
	// Expiration date associated to account
	ExpirationDate string `json:"expirationDate"`
}

// Defines input fields for creating a binding
type CreateBindingInput struct {
	// Name of the binding
	Name string `json:"name"`
	// Principal ID associated with the binding
	PrincipalID uuid.UUID `json:"principalId"`
	// Role ID associated with the binding
	RoleID uuid.UUID `json:"roleId"`
	// Scope reference ID associated with the binding
	ScopeRefID uuid.UUID `json:"scopeRefId"`
	// Version of the binding
	Version string `json:"version"`
}

// Defines input fields for creating a client organization unit
type CreateClientOrganizationUnitInput struct {
	// Description of the client organization unit
	Description *string `json:"description,omitempty"`
	// Name of the client organization unit
	Name string `json:"name"`
	// Parent organization ID
	ParentOrgID uuid.UUID `json:"parentOrgId"`
	// Tenant ID
	TenantID uuid.UUID `json:"tenantId"`
}

// Defines input fields for creating a permission
type CreatePermissionInput struct {
	// Action associated with the permission
	Action string `json:"action"`
	// Unique identifier of the permission
	ID uuid.UUID `json:"id"`
	// Name of the permission
	Name string `json:"name"`
	// Service ID associated with the permission
	ServiceID uuid.UUID `json:"serviceId"`
}

// Defines input fields for creating a role
type CreateRoleInput struct {
	// Assignable scope reference ID
	AssignableScopeRef uuid.UUID `json:"assignableScopeRef"`
	// Description of the role
	Description *string `json:"description,omitempty"`
	// Unique identifier of the role
	ID uuid.UUID `json:"id"`
	// Name of the role
	Name string `json:"name"`
	// Permissions associated with the role
	Permissions []string `json:"permissions"`
	// Type of the role
	RoleType RoleTypeEnum `json:"roleType"`
	// Version of the role
	Version string `json:"version"`
}

// Defines input fields for creating a root
type CreateRootInput struct {
	// Description of the root
	Description *string `json:"description,omitempty"`
	// Name of the root
	Name string `json:"name"`
}

// Defines input fields for creating a tenant
type CreateTenantInput struct {
	// Contact information of the tenant
	ContactInfo *ContactInfoInput `json:"contactInfo,omitempty"`
	// Description of the tenant
	Description *string `json:"description,omitempty"`
	// Unique identifier of the account
	ID uuid.UUID `json:"id"`
	// Name of the tenant
	Name string `json:"name"`
	// Parent organization ID
	ParentID *uuid.UUID `json:"parentId,omitempty"`
}

// Defines input fields for deleting a resource
type DeleteInput struct {
	// Unique identifier of the resource
	ID uuid.UUID `json:"id"`
}

// Represents a Group entity
type Group struct {
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Description of the group
	Description *string `json:"description,omitempty"`
	// Email of the group
	Email string `json:"email"`
	// Unique identifier of the group
	ID uuid.UUID `json:"id"`
	// Members of the group
	Members []*User `json:"members"`
	// Name of the group
	Name string `json:"name"`
	// Tenant associated with the group
	Tenant *Tenant `json:"tenant"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

func (Group) IsData() {}

func (Group) IsPrincipal() {}

// Email of the principal
func (this Group) GetEmail() string { return this.Email }

// Unique identifier of the principal
func (this Group) GetID() uuid.UUID { return this.ID }

// Name of the principal
func (this Group) GetName() string { return this.Name }

// Tenant associated with the principal
func (this Group) GetTenant() *Tenant { return this.Tenant }

func (Group) IsResource() {}

// Timestamp of creation
func (this Group) GetCreatedAt() string { return this.CreatedAt }

// Identifier of the user who created the record
func (this Group) GetCreatedBy() uuid.UUID { return this.CreatedBy }

// Unique identifier of the resource

// Name of the resource

// Timestamp of last update
func (this Group) GetUpdatedAt() string { return this.UpdatedAt }

// Identifier of the user who last updated the record
func (this Group) GetUpdatedBy() uuid.UUID { return this.UpdatedBy }

// Represents a Permission entity
type Permission struct {
	// Action associated with the permission
	Action string `json:"action"`
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Unique identifier of the permission
	ID uuid.UUID `json:"id"`
	// Name of the permission
	Name string `json:"name"`
	// Service ID associated with the permission
	ServiceID string `json:"serviceId"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

func (Permission) IsData() {}

// Define ResponseError for error cases
type ResponseError struct {
	// Error code representing the type of error.
	ErrorCode string `json:"errorCode"`
	// Details about the error.
	ErrorDetails *string `json:"errorDetails,omitempty"`
	// Indicates if the operation was successful.
	IsSuccess bool `json:"isSuccess"`
	// A message providing additional context or information about the operation.
	Message string `json:"message"`
	// A message providing additional context or information about the operation for the logging.
	SystemMessage string `json:"systemMessage"`
}

func (ResponseError) IsOperationResult() {}

func (ResponseError) IsResponse() {}

// Indicates if the operation was successful.
func (this ResponseError) GetIsSuccess() bool { return this.IsSuccess }

// A message providing additional context or information about the operation.
func (this ResponseError) GetMessage() string { return this.Message }

func (ResponseError) IsError() {}

// Error code representing the type of error.
func (this ResponseError) GetErrorCode() string { return this.ErrorCode }

// Details about the error.
func (this ResponseError) GetErrorDetails() *string { return this.ErrorDetails }

// A message providing information about the operation to the user.

// A message providing additional context or information about the operation for the logging.
func (this ResponseError) GetSystemMessage() string { return this.SystemMessage }

// Represents a Role entity
type Role struct {
	// Assignable scope of the role
	AssignableScope Resource `json:"assignableScope"`
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Description of the role
	Description *string `json:"description,omitempty"`
	// Unique identifier of the role
	ID uuid.UUID `json:"id"`
	// Name of the role
	Name string `json:"name"`
	// Permissions associated with the role
	Permissions []*Permission `json:"permissions"`
	// Type of the role
	RoleType RoleTypeEnum `json:"roleType"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
	// Version of the role
	Version string `json:"version"`
}

func (Role) IsData() {}

func (Role) IsResource() {}

// Timestamp of creation
func (this Role) GetCreatedAt() string { return this.CreatedAt }

// Identifier of the user who created the record
func (this Role) GetCreatedBy() uuid.UUID { return this.CreatedBy }

// Unique identifier of the resource
func (this Role) GetID() uuid.UUID { return this.ID }

// Name of the resource
func (this Role) GetName() string { return this.Name }

// Timestamp of last update
func (this Role) GetUpdatedAt() string { return this.UpdatedAt }

// Identifier of the user who last updated the record
func (this Role) GetUpdatedBy() uuid.UUID { return this.UpdatedBy }

// Represents a Root entity
type Root struct {
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Description of the root
	Description *string `json:"description,omitempty"`
	// Unique identifier of the root
	ID uuid.UUID `json:"id"`
	// Name of the root
	Name string `json:"name"`
	// Parent organization
	ParentOrg Organization `json:"parentOrg,omitempty"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

func (Root) IsData() {}

func (Root) IsOrganization() {}

// Timestamp of creation
func (this Root) GetCreatedAt() string { return this.CreatedAt }

// Identifier of the user who created the record
func (this Root) GetCreatedBy() uuid.UUID { return this.CreatedBy }

// Description of the organization
func (this Root) GetDescription() *string { return this.Description }

// Unique identifier of the organization
func (this Root) GetID() uuid.UUID { return this.ID }

// Name of the organization
func (this Root) GetName() string { return this.Name }

// Parent organization
func (this Root) GetParentOrg() Organization { return this.ParentOrg }

// Timestamp of last update
func (this Root) GetUpdatedAt() string { return this.UpdatedAt }

// Identifier of the user who last updated the record
func (this Root) GetUpdatedBy() uuid.UUID { return this.UpdatedBy }

func (Root) IsResource() {}

// Timestamp of creation

// Identifier of the user who created the record

// Unique identifier of the resource

// Name of the resource

// Timestamp of last update

// Identifier of the user who last updated the record

// Success Response for a generic operation
type SuccessResponse struct {
	// The data returned from the operation.
	Data []Data `json:"data,omitempty"`
	// Indicates if the operation was successful.
	IsSuccess bool `json:"isSuccess"`
	// A message providing additional context or information about the operation.
	Message string `json:"message"`
}

func (SuccessResponse) IsOperationResult() {}

func (SuccessResponse) IsResponse() {}

// Indicates if the operation was successful.
func (this SuccessResponse) GetIsSuccess() bool { return this.IsSuccess }

// A message providing additional context or information about the operation.
func (this SuccessResponse) GetMessage() string { return this.Message }

// Represents a Tenant entity
type Tenant struct {
	// Contact information of the tenant
	ContactInfo *ContactInfo `json:"contactInfo,omitempty"`
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Description of the tenant
	Description *string `json:"description,omitempty"`
	// Unique identifier of the tenant
	ID uuid.UUID `json:"id"`
	// Name of the tenant
	Name string `json:"name"`
	// Parent organization
	ParentOrg Organization `json:"parentOrg,omitempty"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

func (Tenant) IsData() {}

func (Tenant) IsOrganization() {}

// Timestamp of creation
func (this Tenant) GetCreatedAt() string { return this.CreatedAt }

// Identifier of the user who created the record
func (this Tenant) GetCreatedBy() uuid.UUID { return this.CreatedBy }

// Description of the organization
func (this Tenant) GetDescription() *string { return this.Description }

// Unique identifier of the organization
func (this Tenant) GetID() uuid.UUID { return this.ID }

// Name of the organization
func (this Tenant) GetName() string { return this.Name }

// Parent organization
func (this Tenant) GetParentOrg() Organization { return this.ParentOrg }

// Timestamp of last update
func (this Tenant) GetUpdatedAt() string { return this.UpdatedAt }

// Identifier of the user who last updated the record
func (this Tenant) GetUpdatedBy() uuid.UUID { return this.UpdatedBy }

func (Tenant) IsResource() {}

// Timestamp of creation

// Identifier of the user who created the record

// Unique identifier of the resource

// Name of the resource

// Timestamp of last update

// Identifier of the user who last updated the record

// Defines input fields for updating an account
type UpdateAccountInput struct {
	// Scope of billing info
	BillingInfo *UpdateBillingInfoInput `json:"billingInfo,omitempty"`
	// Updated description of the account
	Description *string `json:"description,omitempty"`
	// Unique identifier of the account
	ID uuid.UUID `json:"id"`
	// Updated name of the account
	Name *string `json:"name,omitempty"`
	// Associated parent organization
	ParentID *uuid.UUID `json:"parentId,omitempty"`
	// Associated tenant
	TenantID *uuid.UUID `json:"tenantId,omitempty"`
}

// Defines input fields for updating a billing address for an account
type UpdateBillingAddressInput struct {
	// Name of the city associated to billing address
	City *string `json:"city,omitempty"`
	// Name of the country associated to billing address
	Country *string `json:"country,omitempty"`
	// Name of the state associated to billing address
	State *string `json:"state,omitempty"`
	// Name of the street associated to billing address
	Street *string `json:"street,omitempty"`
	// Name of the zipcode associated to billing address
	Zipcode *string `json:"zipcode,omitempty"`
}

// Defines input fields for updating billing info for an account
type UpdateBillingInfoInput struct {
	// Billing Address associated to account
	BillingAddress *UpdateBillingAddressInput `json:"billingAddress,omitempty"`
	// Credit card number associated to account
	CreditCardNumber *string `json:"creditCardNumber,omitempty"`
	// Credit card type associated to account
	CreditCardType *string `json:"creditCardType,omitempty"`
	// CVV associated to account
	Cvv *string `json:"cvv,omitempty"`
	// Expiration date associated to account
	ExpirationDate *string `json:"expirationDate,omitempty"`
}

// Defines input fields for updating a binding
type UpdateBindingInput struct {
	// Unique identifier of the binding
	ID uuid.UUID `json:"id"`
	// Updated name of the binding
	Name string `json:"name"`
	// Updated principal ID associated with the binding
	PrincipalID uuid.UUID `json:"principalId"`
	// Updated role ID associated with the binding
	RoleID uuid.UUID `json:"roleId"`
	// Updated scope reference ID associated with the binding
	ScopeRefID uuid.UUID `json:"scopeRefId"`
	// Updated version of the binding
	Version string `json:"version"`
}

// Defines input fields for updating a client organization unit
type UpdateClientOrganizationUnitInput struct {
	// Updated description of the client organization unit
	Description *string `json:"description,omitempty"`
	// Unique identifier of the client organization unit
	ID uuid.UUID `json:"id"`
	// Updated name of the client organization unit
	Name *string `json:"name,omitempty"`
	// Updated parent organization ID
	ParentOrgID *uuid.UUID `json:"parentOrgId,omitempty"`
	// Updated tenant ID
	TenantID *uuid.UUID `json:"tenantId,omitempty"`
}

// Defines input fields for updating a permission
type UpdatePermissionInput struct {
	// Updated action associated with the permission
	Action string `json:"action"`
	// Unique identifier of the permission
	ID uuid.UUID `json:"id"`
	// Updated name of the permission
	Name string `json:"name"`
	// Updated service ID associated with the permission
	ServiceID uuid.UUID `json:"serviceId"`
}

// Defines input fields for updating a role
type UpdateRoleInput struct {
	// Updated assignable scope reference ID
	AssignableScopeRef uuid.UUID `json:"assignableScopeRef"`
	// Updated description of the role
	Description *string `json:"description,omitempty"`
	// Unique identifier of the role
	ID uuid.UUID `json:"id"`
	// Updated name of the role
	Name string `json:"name"`
	// Updated permissions associated with the role
	Permissions []string `json:"permissions"`
	// Updated type of the role
	RoleType RoleTypeEnum `json:"roleType"`
	// Updated version of the role
	Version string `json:"version"`
}

// Defines input fields for updating a root
type UpdateRootInput struct {
	// Updated description of the root
	Description *string `json:"description,omitempty"`
	// Unique identifier of the root
	ID uuid.UUID `json:"id"`
	// Updated name of the root
	Name *string `json:"name,omitempty"`
}

// Defines input fields for updating a tenant
type UpdateTenantInput struct {
	// Updated contact information of the tenant
	ContactInfo *ContactInfoInput `json:"contactInfo,omitempty"`
	// Updated description of the tenant
	Description *string `json:"description,omitempty"`
	// Unique identifier of the tenant
	ID uuid.UUID `json:"id"`
	// Updated name of the tenant
	Name *string `json:"name,omitempty"`
	// Updated parent organization ID
	ParentID *uuid.UUID `json:"parentId,omitempty"`
}

// Represents a User entity
type User struct {
	// Timestamp of creation
	CreatedAt string `json:"createdAt"`
	// Identifier of the user who created the record
	CreatedBy uuid.UUID `json:"createdBy"`
	// Email of the user
	Email string `json:"email"`
	// First name of the user
	FirstName string `json:"firstName"`
	// Unique identifier of the user
	ID uuid.UUID `json:"id"`
	// Last name of the user
	LastName string `json:"lastName"`
	// Name of the user
	Name string `json:"name"`
	// Tenant associated with the user
	Tenant *Tenant `json:"tenant"`
	// Timestamp of last update
	UpdatedAt string `json:"updatedAt"`
	// Identifier of the user who last updated the record
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

func (User) IsData() {}

func (User) IsPrincipal() {}

// Email of the principal
func (this User) GetEmail() string { return this.Email }

// Unique identifier of the principal
func (this User) GetID() uuid.UUID { return this.ID }

// Name of the principal
func (this User) GetName() string { return this.Name }

// Tenant associated with the principal
func (this User) GetTenant() *Tenant { return this.Tenant }

func (User) IsResource() {}

// Timestamp of creation
func (this User) GetCreatedAt() string { return this.CreatedAt }

// Identifier of the user who created the record
func (this User) GetCreatedBy() uuid.UUID { return this.CreatedBy }

// Unique identifier of the resource

// Name of the resource

// Timestamp of last update
func (this User) GetUpdatedAt() string { return this.UpdatedAt }

// Identifier of the user who last updated the record
func (this User) GetUpdatedBy() uuid.UUID { return this.UpdatedBy }

// Defines the role type enumeration
type RoleTypeEnum string

const (
	// Custom role type
	RoleTypeEnumCustom RoleTypeEnum = "CUSTOM"
	// Default role type
	RoleTypeEnumDefault RoleTypeEnum = "DEFAULT"
)

var AllRoleTypeEnum = []RoleTypeEnum{
	RoleTypeEnumCustom,
	RoleTypeEnumDefault,
}

func (e RoleTypeEnum) IsValid() bool {
	switch e {
	case RoleTypeEnumCustom, RoleTypeEnumDefault:
		return true
	}
	return false
}

func (e RoleTypeEnum) String() string {
	return string(e)
}

func (e *RoleTypeEnum) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RoleTypeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RoleTypeEnum", str)
	}
	return nil
}

func (e RoleTypeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
