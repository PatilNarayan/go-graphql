package controller

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"
	"fmt"
	"go_graphql/graph/model"
)

// CreateTenant is the resolver for the createTenant field.
func (r *mutationResolver) CreateTenant(ctx context.Context, input model.TenantInput) (*model.Tenant, error) {
	// Call the PermitClient to create a tenant
	attributes := make(map[string]interface{})
	tenant, err := r.PermitClient.CreateTenant(&input.Name, &input.Name, attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %v", err)
	}

	// Map the response from PermitClient into your model.Tenant object
	createdTenant := &model.Tenant{
		ID:   tenant.Id, // Assuming tenant has an ID field
		Name: tenant.Name,
	}

	return createdTenant, nil
}

// UpdateTenant is the resolver for the updateTenant field.
func (r *mutationResolver) UpdateTenant(ctx context.Context, id string, input model.TenantInput) (*model.Tenant, error) {
	// Call PermitClient to update a tenant by ID
	attributes := make(map[string]interface{})
	updatedTenant, err := r.PermitClient.UpdateTenant(id, &input.Name, &input.Name, attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to update tenant: %v", err)
	}

	// Map response into the model
	updated := &model.Tenant{
		ID:   updatedTenant.Id,
		Name: updatedTenant.Name,
	}

	return updated, nil
}

// DeleteTenant is the resolver for the deleteTenant field.
func (r *mutationResolver) DeleteTenant(ctx context.Context, id string) (bool, error) {
	// Call PermitClient to delete the tenant
	err := r.PermitClient.DeleteTenant(id)
	if err != nil {
		return false, fmt.Errorf("failed to delete tenant: %v", err)
	}

	// Return success
	return true, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	// Call the PermitClient to create a user
	attributes := make(map[string]interface{})
	rolls := []map[string]interface{}{}
	user, err := r.PermitClient.CreateUser(input.Email, input.FirstName, input.LastName, attributes, rolls)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Map response from PermitClient into the model.User object
	createdUser := &model.User{
		Key:       user.Id, // Assuming user has ID field
		Email:     *user.Email,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
	}

	return createdUser, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, key string, input model.UserInput) (*model.User, error) {
	// Call PermitClient to update the user by key
	attributes := make(map[string]interface{})
	rolls := []map[string]interface{}{}
	updatedUser, err := r.PermitClient.UpdateUser(key, &input.Email, &input.FirstName, &input.LastName, attributes, rolls)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Map the response into the model.User object

	updated := &model.User{
		Key:       updatedUser.Id,
		Email:     *updatedUser.Email,
		FirstName: *updatedUser.FirstName,
		LastName:  *updatedUser.LastName,
	}

	return updated, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, key string) (bool, error) {
	// Call PermitClient to delete the user by key
	err := r.PermitClient.DeleteUser(key)
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %v", err)
	}

	// Return success
	return true, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
