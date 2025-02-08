package gql

import (
	"go_graphql/gql/generated"
	"go_graphql/internal/accounts"
	"go_graphql/internal/permissions"
	"go_graphql/internal/role"
	"go_graphql/internal/tenants"

	"gorm.io/gorm"
)

// Resolver holds references to the DB and acts as a central resolver
type Resolver struct {
	DB *gorm.DB
}

// Query returns the root query resolvers, delegating to feature-based resolverssss
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{

		TenantQueryResolver:  &tenants.TenantQueryResolver{DB: r.DB},
		AccountQueryResolver: &accounts.AccountQueryResolver{DB: r.DB},
		// ClientOrganizationUnitQueryResolver: &clientorganizationunit.ClientOrganizationUnitQueryResolver{DB: r.DB},
		RoleQueryResolver:       &role.RoleQueryResolver{DB: r.DB},
		PermissionQueryResolver: &permissions.PermissionQueryResolver{DB: r.DB},
		// ResourceQueryResolver:               &resources.ResourceQueryResolver{DB: r.DB},
	}
}

// Mutation returns the root mutation resolvers, delegating to feature-based resolvers
func (r *Resolver) Mutation() generated.MutationResolver {
	// permitClient := permit.NewPermitClient()
	return &mutationResolver{
		TenantMutationResolver: &tenants.TenantMutationResolver{DB: r.DB},
		// AccountMutationResolver: &accounts.AccountMutationResolver{DB: r.DB},
		// ClientOrganizationUnitMutationResolver: &clientorganizationunit.ClientOrganizationUnitMutationResolver{r.DB},
		RoleMutationResolver:       &role.RoleMutationResolver{DB: r.DB},
		PermissionMutationResolver: &permissions.PermissionMutationResolver{DB: r.DB},
		// ResourceInstanceMutationResolver: &resourceInstances.ResourceInstanceMutationResolver{DB: r.DB},
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

// Account resolves fields for the Account type
func (r *Resolver) Account() generated.AccountResolver {
	return &accounts.AccountFieldResolver{DB: r.DB}
}

type AccountResolver struct{ *Resolver }

// Root resolvers for Query and Mutation
type queryResolver struct {
	*tenants.TenantQueryResolver
	*accounts.AccountQueryResolver
	*role.RoleQueryResolver
	// *clientorganizationunit.ClientOrganizationUnitQueryResolver
	*permissions.PermissionQueryResolver
	// *resources.ResourceQueryResolver
}

type mutationResolver struct {
	*tenants.TenantMutationResolver
	*accounts.AccountMutationResolver
	// *clientorganizationunit.ClientOrganizationUnitMutationResolver
	*role.RoleMutationResolver
	*permissions.PermissionMutationResolver
	// *resourceInstances.ResourceInstanceMutationResolver
}
