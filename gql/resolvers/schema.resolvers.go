package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"
	"fmt"
	"go_graphql/gql/generated"
	"go_graphql/gql/models"

	"github.com/google/uuid"
)

// CreateTenant is the resolver for the createTenant field.
func (r *mutationResolver) CreateTenant(ctx context.Context, input models.CreateTenantInput) (*models.Tenant, error) {
	panic(fmt.Errorf("not implemented: CreateTenant - createTenant"))
}

// UpdateTenant is the resolver for the updateTenant field.
func (r *mutationResolver) UpdateTenant(ctx context.Context, input models.UpdateTenantInput) (*models.Tenant, error) {
	panic(fmt.Errorf("not implemented: UpdateTenant - updateTenant"))
}

// DeleteTenant is the resolver for the deleteTenant field.
func (r *mutationResolver) DeleteTenant(ctx context.Context, id uuid.UUID) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteTenant - deleteTenant"))
}

// CreateRole is the resolver for the createRole field.
func (r *mutationResolver) CreateRole(ctx context.Context, input models.CreateRoleInput) (*models.Role, error) {
	panic(fmt.Errorf("not implemented: CreateRole - createRole"))
}

// UpdateRole is the resolver for the updateRole field.
func (r *mutationResolver) UpdateRole(ctx context.Context, input models.UpdateRoleInput) (*models.Role, error) {
	panic(fmt.Errorf("not implemented: UpdateRole - updateRole"))
}

// DeleteRole is the resolver for the deleteRole field.
func (r *mutationResolver) DeleteRole(ctx context.Context, id uuid.UUID) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteRole - deleteRole"))
}

// CreatePermission is the resolver for the createPermission field.
func (r *mutationResolver) CreatePermission(ctx context.Context, input *models.CreatePermission) (*models.Permission, error) {
	panic(fmt.Errorf("not implemented: CreatePermission - createPermission"))
}

// DeletePermission is the resolver for the deletePermission field.
func (r *mutationResolver) DeletePermission(ctx context.Context, id uuid.UUID) (bool, error) {
	panic(fmt.Errorf("not implemented: DeletePermission - deletePermission"))
}

// UpdatePermission is the resolver for the updatePermission field.
func (r *mutationResolver) UpdatePermission(ctx context.Context, input *models.UpdatePermission) (*models.Permission, error) {
	panic(fmt.Errorf("not implemented: UpdatePermission - updatePermission"))
}

// GetTenant is the resolver for the getTenant field.
func (r *queryResolver) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	panic(fmt.Errorf("not implemented: GetTenant - getTenant"))
}

// AllTenants is the resolver for the allTenants field.
func (r *queryResolver) AllTenants(ctx context.Context) ([]*models.Tenant, error) {
	panic(fmt.Errorf("not implemented: AllTenants - allTenants"))
}

// GetRole is the resolver for the getRole field.
func (r *queryResolver) GetRole(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	panic(fmt.Errorf("not implemented: GetRole - getRole"))
}

// AllRoles is the resolver for the allRoles field.
func (r *queryResolver) AllRoles(ctx context.Context) ([]*models.Role, error) {
	panic(fmt.Errorf("not implemented: AllRoles - allRoles"))
}

// GetAllRolesForAssignableScopeRef is the resolver for the getAllRolesForAssignableScopeRef field.
func (r *queryResolver) GetAllRolesForAssignableScopeRef(ctx context.Context, id uuid.UUID) ([]*models.Role, error) {
	panic(fmt.Errorf("not implemented: GetAllRolesForAssignableScopeRef - getAllRolesForAssignableScopeRef"))
}

// GetAllPermissions is the resolver for the getAllPermissions field.
func (r *queryResolver) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	panic(fmt.Errorf("not implemented: GetAllPermissions - getAllPermissions"))
}

// GetPermission is the resolver for the getPermission field.
func (r *queryResolver) GetPermission(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	panic(fmt.Errorf("not implemented: GetPermission - getPermission"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
