package permit

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/models"
	"github.com/permitio/permit-golang/pkg/permit"
	"github.com/rs/xid"
)

type PermitClient struct {
	client *permit.Client
}

// NewPermitClient initializes the PermitClient with API token
func NewPermitClient() (*PermitClient, error) {
	// Create config with provided API token
	permitConfig := config.NewConfigBuilder(os.Getenv("PERMIT_TOKEN")).Build()

	// Initialize the Permit client
	client := permit.New(permitConfig)

	return &PermitClient{client: client}, nil
}

// CreateTenant creates a new tenant on Permit.io
func (pc *PermitClient) CreateTenant(name, description *string, attributes map[string]interface{}) (*models.TenantRead, error) {
	xkey := xid.New().String()
	tenantCreate := models.NewTenantCreate(xkey, *name)
	tenantCreate.SetName(*name)
	if description != nil {
		tenantCreate.SetDescription(*description)
	}
	if attributes != nil {
		tenantCreate.SetAttributes(attributes)
	}

	tenant, err := pc.client.Api.Tenants.Create(context.TODO(), *tenantCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %v", err)
	}
	log.Printf("User created successfully: %v", tenant)

	return tenant, nil
}

// UpdateTenant updates an existing tenant on Permit.io
func (pc *PermitClient) UpdateTenant(tenantKey string, name, description *string, attributes map[string]interface{}) (*models.TenantRead, error) {
	// Create an update model for tenant
	tenantUpdate := models.TenantUpdate{}

	// Set the fields to be updated if provided
	if name != nil {
		tenantUpdate.SetName(*name)
	}
	if description != nil {
		tenantUpdate.SetDescription(*description)
	}
	if attributes != nil {
		tenantUpdate.SetAttributes(attributes)
	}

	// Update tenant
	updatedTenant, err := pc.client.Api.Tenants.Update(context.TODO(), tenantKey, tenantUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to update tenant: %v", err)
	}

	log.Printf("Tenant updated successfully: %v", updatedTenant)
	return updatedTenant, nil
}

// DeleteTenant deletes a tenant from Permit.io
func (pc *PermitClient) DeleteTenant(tenantKey string) error {
	// Delete tenant by tenant key
	err := pc.client.Api.Tenants.Delete(context.TODO(), tenantKey)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %v", err)
	}

	log.Print("Tenant deleted successfully.")
	return nil
}

// CreateUser creates a new user in Permit.io
func (pc *PermitClient) CreateUser(email, firstName, lastName string, attributes map[string]interface{}, roleAssignments []map[string]interface{}) (*models.UserRead, error) {
	userKey := xid.New().String()
	userInput := models.NewUserCreate(userKey)
	userInput.Email = &email
	userInput.FirstName = &firstName
	userInput.LastName = &lastName
	userInput.SetAttributes(attributes)

	newUser, err := pc.client.Api.Users.Create(context.TODO(), *userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	log.Printf("User created successfully: %v", newUser)
	return newUser, nil
}

// UpdateUser updates an existing user in Permit.io
func (pc *PermitClient) UpdateUser(userKey string, email, firstName, lastName *string, attributes map[string]interface{}, roleAssignments []map[string]interface{}) (*models.UserRead, error) {
	// Retrieve the current user first, if needed, or assume you already have the `userKey`
	userInput := models.UserUpdate{}
	if email != nil {
		userInput.Email = email
	}
	if firstName != nil {
		userInput.FirstName = firstName
	}
	if lastName != nil {
		userInput.LastName = lastName
	}

	if attributes != nil {
		userInput.SetAttributes(attributes)
	}

	// Update user
	updatedUser, err := pc.client.Api.Users.Update(context.TODO(), userKey, userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	log.Printf("User updated successfully: %v", updatedUser)
	return updatedUser, nil
}

// DeleteUser deletes a user from Permit.io
func (pc *PermitClient) DeleteUser(userKey string) error {
	// Delete user by user key
	err := pc.client.Api.Users.Delete(context.TODO(), userKey)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	log.Print("User deleted successfully: ")
	return nil
}
