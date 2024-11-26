package gql

import (
	"context"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"iam_services_main_v1/gql/generated"
	"iam_services_main_v1/internal/groups"
	"iam_services_main_v1/internal/organizations"
	"iam_services_main_v1/internal/tenants"

	"gorm.io/gorm"
)

// Resolver holds references to the DB and acts as a central resolver
type Resolver struct {
	DB *gorm.DB
}

// Query returns the root query resolvers, delegating to feature-based resolvers
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{
		OrganizationQueryResolver: &organizations.OrganizationQueryResolver{DB: r.DB},
		TenantQueryResolver:       &tenants.TenantQueryResolver{DB: r.DB},
		GroupQueryResolver:        &groups.GroupQueryResolver{DB: r.DB},
	}
}

// Mutation returns the root mutation resolvers, delegating to feature-based resolvers
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{
		OrganizationMutationResolver: &organizations.OrganizationMutationResolver{DB: r.DB},
		TenantMutationResolver:       &tenants.TenantMutationResolver{DB: r.DB},
		GroupMutationResolver:        &groups.GroupMutationResolver{DB: r.DB},
	}
}

// // Organization resolves fields for the Organization type
// func (r *Resolver) Organization() generated.OrganizationResolver {
// return &organizations.OrganizationFieldResolver{DB: r.DB}
// }

// Organization resolves fields for the Organization type
func (r *Resolver) Tenant() generated.TenantResolver {
	return &tenants.TenantFieldResolver{DB: r.DB}
}

// Organization resolves fields for the Organization type
func (r *Resolver) Group() generated.GroupResolver {
	return &groups.GroupFieldResolver{DB: r.DB}
}

// Root resolvers for Query and Mutation
type queryResolver struct {
	*organizations.OrganizationQueryResolver
	*tenants.TenantQueryResolver
	*groups.GroupQueryResolver
}

// GetGroup implements generated.QueryResolver.
// Subtle: this method shadows the method (*GroupQueryResolver).GetGroup of queryResolver.GroupQueryResolver.
func (q *queryResolver) GetGroup(ctx context.Context, id *int) (*dto.GroupEntity, error) {
	panic("unimplemented")
}

// GetOrganization implements generated.QueryResolver.
// Subtle: this method shadows the method (*OrganizationQueryResolver).GetOrganization of queryResolver.OrganizationQueryResolver.
func (q *queryResolver) GetOrganization(ctx context.Context, id *int) (dto.Organization, error) {
	panic("unimplemented")
}

// GetTenant implements generated.QueryResolver.
// Subtle: this method shadows the method (*TenantQueryResolver).GetTenant of queryResolver.TenantQueryResolver.
func (q *queryResolver) GetTenant(ctx context.Context, id *int) (*dto.Tenant, error) {
	panic("unimplemented")
}

// Groups implements generated.QueryResolver.
// Subtle: this method shadows the method (*GroupQueryResolver).Groups of queryResolver.GroupQueryResolver.
func (q *queryResolver) Groups(ctx context.Context) ([]*dto.GroupEntity, error) {
	panic("unimplemented")
}

// Organizations implements generated.QueryResolver.
// Subtle: this method shadows the method (*OrganizationQueryResolver).Organizations of queryResolver.OrganizationQueryResolver.
func (q *queryResolver) Organizations(ctx context.Context) ([]dto.Organization, error) {
	panic("unimplemented")
}

// Tenants implements generated.QueryResolver.
// Subtle: this method shadows the method (*TenantQueryResolver).Tenants of queryResolver.TenantQueryResolver.
func (q *queryResolver) Tenants(ctx context.Context) ([]*dto.Tenant, error) {
	panic("unimplemented")
}

type mutationResolver struct {
	*organizations.OrganizationMutationResolver
	*tenants.TenantMutationResolver
	*groups.GroupMutationResolver
}

// CreateGroup implements generated.MutationResolver.
// Subtle: this method shadows the method (*GroupMutationResolver).CreateGroup of mutationResolver.GroupMutationResolver.
func (m *mutationResolver) CreateGroup(ctx context.Context, input models.GroupInput) (*dto.GroupEntity, error) {
	panic("unimplemented")
}

// CreateOrganization implements generated.MutationResolver.
// Subtle: this method shadows the method (*OrganizationMutationResolver).CreateOrganization of mutationResolver.OrganizationMutationResolver.
func (m *mutationResolver) CreateOrganization(ctx context.Context, name string) (dto.Organization, error) {
	panic("unimplemented")
}

// CreateTenant implements generated.MutationResolver.
// Subtle: this method shadows the method (*TenantMutationResolver).CreateTenant of mutationResolver.TenantMutationResolver.
func (m *mutationResolver) CreateTenant(ctx context.Context, input models.TenantInput) (*dto.Tenant, error) {
	panic("unimplemented")
}

// DeleteGroup implements generated.MutationResolver.
// Subtle: this method shadows the method (*GroupMutationResolver).DeleteGroup of mutationResolver.GroupMutationResolver.
func (m *mutationResolver) DeleteGroup(ctx context.Context, id int) (bool, error) {
	panic("unimplemented")
}

// DeleteTenant implements generated.MutationResolver.
// Subtle: this method shadows the method (*TenantMutationResolver).DeleteTenant of mutationResolver.TenantMutationResolver.
func (m *mutationResolver) DeleteTenant(ctx context.Context, id int) (bool, error) {
	panic("unimplemented")
}

// UpdateGroup implements generated.MutationResolver.
// Subtle: this method shadows the method (*GroupMutationResolver).UpdateGroup of mutationResolver.GroupMutationResolver.
func (m *mutationResolver) UpdateGroup(ctx context.Context, id int, input models.GroupInput) (*dto.GroupEntity, error) {
	panic("unimplemented")
}

// UpdateTenant implements generated.MutationResolver.
// Subtle: this method shadows the method (*TenantMutationResolver).UpdateTenant of mutationResolver.TenantMutationResolver.
func (m *mutationResolver) UpdateTenant(ctx context.Context, id int, input models.TenantInput) (*dto.Tenant, error) {
	panic("unimplemented")
}

// Organization resolves fields for the Organization type
func (r *Resolver) Group() generated.GroupResolver {
	return &groups.GroupFieldResolver{DB: r.DB}
}

// Root resolvers for Query and Mutation
type queryResolver struct {
	*organizations.OrganizationQueryResolver
	*tenants.TenantQueryResolver
	*groups.GroupQueryResolver
}

// GetGroup implements generated.QueryResolver.
// Subtle: this method shadows the method (*GroupQueryResolver).GetGroup of queryResolver.GroupQueryResolver.
func (q *queryResolver) GetGroup(ctx context.Context, id *int) (*dto.GroupEntity, error) {
	panic("unimplemented")
}

// GetOrganization implements generated.QueryResolver.
// Subtle: this method shadows the method (*OrganizationQueryResolver).GetOrganization of queryResolver.OrganizationQueryResolver.
func (q *queryResolver) GetOrganization(ctx context.Context, id *int) (dto.Organization, error) {
	panic("unimplemented")
}

// GetTenant implements generated.QueryResolver.
// Subtle: this method shadows the method (*TenantQueryResolver).GetTenant of queryResolver.TenantQueryResolver.
func (q *queryResolver) GetTenant(ctx context.Context, id *int) (*dto.Tenant, error) {
	panic("unimplemented")
}

// Groups implements generated.QueryResolver.
// Subtle: this method shadows the method (*GroupQueryResolver).Groups of queryResolver.GroupQueryResolver.
func (q *queryResolver) Groups(ctx context.Context) ([]*dto.GroupEntity, error) {
	panic("unimplemented")
}

// Organizations implements generated.QueryResolver.
// Subtle: this method shadows the method (*OrganizationQueryResolver).Organizations of queryResolver.OrganizationQueryResolver.
func (q *queryResolver) Organizations(ctx context.Context) ([]dto.Organization, error) {
	panic("unimplemented")
}

// Tenants implements generated.QueryResolver.
// Subtle: this method shadows the method (*TenantQueryResolver).Tenants of queryResolver.TenantQueryResolver.
func (q *queryResolver) Tenants(ctx context.Context) ([]*dto.Tenant, error) {
	panic("unimplemented")
}

type mutationResolver struct {
	*organizations.OrganizationMutationResolver
	*tenants.TenantMutationResolver
	*groups.GroupMutationResolver
}

// CreateGroup implements generated.MutationResolver.
// Subtle: this method shadows the method (*GroupMutationResolver).CreateGroup of mutationResolver.GroupMutationResolver.
func (m *mutationResolver) CreateGroup(ctx context.Context, input models.GroupInput) (*dto.GroupEntity, error) {
	panic("unimplemented")
}

// CreateOrganization implements generated.MutationResolver.
// Subtle: this method shadows the method (*OrganizationMutationResolver).CreateOrganization of mutationResolver.OrganizationMutationResolver.
func (m *mutationResolver) CreateOrganization(ctx context.Context, name string) (dto.Organization, error) {
	panic("unimplemented")
}

// CreateTenant implements generated.MutationResolver.
// Subtle: this method shadows the method (*TenantMutationResolver).CreateTenant of mutationResolver.TenantMutationResolver.
func (m *mutationResolver) CreateTenant(ctx context.Context, input models.TenantInput) (*dto.Tenant, error) {
	panic("unimplemented")
}

// DeleteGroup implements generated.MutationResolver.
// Subtle: this method shadows the method (*GroupMutationResolver).DeleteGroup of mutationResolver.GroupMutationResolver.
func (m *mutationResolver) DeleteGroup(ctx context.Context, id int) (bool, error) {
	panic("unimplemented")
}

// DeleteTenant implements generated.MutationResolver.
// Subtle: this method shadows the method (*TenantMutationResolver).DeleteTenant of mutationResolver.TenantMutationResolver.
func (m *mutationResolver) DeleteTenant(ctx context.Context, id int) (bool, error) {
	panic("unimplemented")
}

// UpdateGroup implements generated.MutationResolver.
// Subtle: this method shadows the method (*GroupMutationResolver).UpdateGroup of mutationResolver.GroupMutationResolver.
func (m *mutationResolver) UpdateGroup(ctx context.Context, id int, input models.GroupInput) (*dto.GroupEntity, error) {
	panic("unimplemented")
}

// UpdateTenant implements generated.MutationResolver.
// Subtle: this method shadows the method (*TenantMutationResolver).UpdateTenant of mutationResolver.TenantMutationResolver.
func (m *mutationResolver) UpdateTenant(ctx context.Context, id int, input models.TenantInput) (*dto.Tenant, error) {
	panic("unimplemented")
}
