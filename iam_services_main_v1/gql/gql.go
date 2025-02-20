package gql

import (
	"iam_services_main_v1/gormlogger"
	"iam_services_main_v1/gql/generated"
	"iam_services_main_v1/internal/accounts"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/internal/roles"
	"iam_services_main_v1/internal/tenants"

	"gorm.io/gorm"
)

// Resolver holds references to the DB and acts as a central resolver
type Resolver struct {
	DB     *gorm.DB
	PC     *permit.PermitClient
	Logger *gormlogger.GORMLogger
}

// Query returns the root query resolvers, delegating to feature-based resolvers
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{
		TenantQueryResolver: &tenants.TenantQueryResolver{DB: r.DB, PC: r.PC},
		// AccountQueryResolver:                &accounts.AccountQueryResolver{DB: r.DB, PC: r.PC},
		// ClientOrganizationUnitQueryResolver: &clientorganizationunits.ClientOrganizationUnitQueryResolver{DB: r.DB},
		RoleQueryResolver: &roles.RoleQueryResolver{DB: r.DB},
		// PermissionQueryResolver:             &permissions.PermissionQueryResolver{DB: r.DB},
		// BindingsQueryResolver:               &bindings.BindingsQueryResolver{DB: r.DB},
		// ResourceQueryResolver:               &resources.ResourceQueryResolver{DB: r.DB},
		// GroupQueryResolver:                  &groups.GroupQueryResolver{DB: r.DB},
		// OrganizationQueryResolver:           &organizations.OrganizationQueryResolver{DB: r.DB},
		// RootQueryResolver:                   &root.RootQueryResolver{DB: r.DB},
	}
}

// Mutation returns the root mutation resolvers, delegating to feature-based resolvers
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{

		TenantMutationResolver: &tenants.TenantMutationResolver{DB: r.DB, PermitClient: r.PC},
		// AccountMutationResolver:                &accounts.AccountMutationResolver{DB: r.DB, PC: r.PC},
		// ClientOrganizationUnitMutationResolver: &clientorganizationunits.ClientOrganizationUnitMutationResolver{r.DB},
		RoleMutationResolver: &roles.RoleMutationResolver{DB: r.DB},
		// PermissionMutationResolver:             &permissions.PermissionMutationResolver{DB: r.DB},
		// BindingsMutationResolver:               &bindings.BindingsMutationResolver{DB: r.DB},
		// RootMutationResolver:                   &root.RootMutationResolver{DB: r.DB},
	}
}

// Account resolves fields for the Account type
func (r *Resolver) Account() generated.AccountResolver {
	return &accounts.AccountFieldResolver{DB: r.DB}
}

type AccountResolver struct{ *Resolver }

// Root resolvers for Query and Mutation
type queryResolver struct {
	*tenants.TenantQueryResolver
	// *accounts.AccountQueryResolver
	*roles.RoleQueryResolver
	// *clientorganizationunits.ClientOrganizationUnitQueryResolver
	// *permissions.PermissionQueryResolver
	// *bindings.BindingsQueryResolver
	// *resources.ResourceQueryResolver
	// *groups.GroupQueryResolver
	// *organizations.OrganizationQueryResolver
	// *root.RootQueryResolver
}

type mutationResolver struct {
	*tenants.TenantMutationResolver
	// *accounts.AccountMutationResolver
	// *clientorganizationunits.ClientOrganizationUnitMutationResolver
	*roles.RoleMutationResolver
	// *permissions.PermissionMutationResolver
	// *bindings.BindingsMutationResolver
	// *root.RootMutationResolver
}
