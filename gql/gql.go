package gql

import (
	"go_graphql/gql/generated"
	clientorganizationunit "go_graphql/internal/clientOrganizationUnit"
	"go_graphql/internal/groups"
	"go_graphql/internal/organizations"
	"go_graphql/internal/role"
	"go_graphql/internal/tenants"

	"gorm.io/gorm"
)

// Resolver holds references to the DB and acts as a central resolver
type Resolver struct {
	DB *gorm.DB
}

// Query returns the root query resolvers, delegating to feature-based resolvers
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{
		OrganizationQueryResolver:           &organizations.OrganizationQueryResolver{DB: r.DB},
		TenantQueryResolver:                 &tenants.TenantQueryResolver{DB: r.DB},
		GroupQueryResolver:                  &groups.GroupQueryResolver{DB: r.DB},
		ClientOrganizationUnitQueryResolver: &clientorganizationunit.ClientOrganizationUnitQueryResolver{DB: r.DB},
		RoleQueryResolver:                   &role.RoleQueryResolver{DB: r.DB},
	}
}

// Mutation returns the root mutation resolvers, delegating to feature-based resolvers
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{
		OrganizationMutationResolver: &organizations.OrganizationMutationResolver{DB: r.DB},
		TenantMutationResolver:       &tenants.TenantMutationResolver{DB: r.DB},
		GroupMutationResolver:        &groups.GroupMutationResolver{DB: r.DB},
		RoleMutationResolver:         &role.RoleMutationResolver{DB: r.DB},
	}
}

// // Organization resolves fields for the Organization type
// func (r *Resolver) Organization() generated.OrganizationResolver {
// return &organizations.OrganizationFieldResolver{DB: r.DB}
// }

// Organization resolves fields for the Organization type
// func (r *Resolver) Tenant() generated.TenantResolver {
// 	return &tenants.TenantFieldResolver{DB: r.DB}
// }

// Organization resolves fields for the Organization type
// func (r *Resolver) Group() generated.GroupResolver {
// 	return &groups.GroupFieldResolver{DB: r.DB}
// }

// Root resolvers for Query and Mutation
type queryResolver struct {
	*organizations.OrganizationQueryResolver
	*tenants.TenantQueryResolver
	*groups.GroupQueryResolver
	*role.RoleQueryResolver
	*clientorganizationunit.ClientOrganizationUnitQueryResolver
}

type mutationResolver struct {
	*organizations.OrganizationMutationResolver
	*tenants.TenantMutationResolver
	*groups.GroupMutationResolver
	*role.RoleMutationResolver
}
