package gql

import (
	"iam_services_main_v1/gql/generated"
	"iam_services_main_v1/internal/accounts"
	"iam_services_main_v1/internal/bindings"
	clientorganizationunit "iam_services_main_v1/internal/clientOrganizationUnit"
	permission "iam_services_main_v1/internal/permissions"
	"iam_services_main_v1/internal/permit"
	role "iam_services_main_v1/internal/roles"
	"iam_services_main_v1/internal/tenants"

	"gorm.io/gorm"
)

// Resolver holds references to the DB and acts as a central resolver
type Resolver struct {
	DB     *gorm.DB
	Permit *permit.PermitClient
}

// Query returns the root query resolvers, delegating to feature-based resolverssss
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{
		TenantQueryResolver:                 &tenants.TenantQueryResolver{DB: r.DB},
		AccountQueryResolver:                &accounts.AccountQueryResolver{DB: r.DB},
		ClientOrganizationUnitQueryResolver: &clientorganizationunit.ClientOrganizationUnitQueryResolver{DB: r.DB},
		RoleQueryResolver:                   &role.RoleQueryResolver{DB: r.DB},
		PermissionQueryResolver:             &permission.PermissionQueryResolver{DB: r.DB, Permit: r.Permit},
		BindingsQueryResolver:               &bindings.BindingsQueryResolver{DB: r.DB},
	}
}

// Mutation returns the root mutation resolvers, delegating to feature-based resolvers
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{
		TenantMutationResolver:                 &tenants.TenantMutationResolver{DB: r.DB},
		AccountMutationResolver:                &accounts.AccountMutationResolver{DB: r.DB},
		ClientOrganizationUnitMutationResolver: &clientorganizationunit.ClientOrganizationUnitMutationResolver{DB: r.DB},
		RoleMutationResolver:                   &role.RoleMutationResolver{DB: r.DB},
		PermissionMutationResolver:             &permission.PermissionMutationResolver{DB: r.DB, Permit: r.Permit},
		BindingsMutationResolver:               &bindings.BindingsMutationResolver{DB: r.DB},
	}
}

func (r *Resolver) ClientOrganizationUnit() generated.ClientOrganizationUnitResolver {
	return &clientorganizationunit.ClientOrganizationUnitFieldResolver{DB: r.DB}
}

func (r *Resolver) Tenant() generated.TenantResolver {
	return &tenants.TenantFieldResolver{DB: r.DB}
}

// Organization resolves fields for the Organization type
// func (r *Resolver) Group() generated.GroupResolver {
// 	return &groups.GroupFieldResolver{DB: r.DB}
// }

// Account resolves fields for the Account type
func (r *Resolver) Account() generated.AccountResolver {
	return &accounts.AccountFieldResolver{DB: r.DB}
}

func (r *Resolver) Binding() generated.BindingResolver {
	// Implement the Binding resolver logic here
	return &bindings.BindingsFieldResolver{DB: r.DB}
}

type AccountResolver struct{ *Resolver }

// Root resolvers for Query and Mutation
type queryResolver struct {
	*tenants.TenantQueryResolver
	*accounts.AccountQueryResolver
	*role.RoleQueryResolver
	*clientorganizationunit.ClientOrganizationUnitQueryResolver
	*permission.PermissionQueryResolver
	*bindings.BindingsQueryResolver
}

type mutationResolver struct {
	*tenants.TenantMutationResolver
	*accounts.AccountMutationResolver
	*clientorganizationunit.ClientOrganizationUnitMutationResolver
	*role.RoleMutationResolver
	*permission.PermissionMutationResolver
	*bindings.BindingsMutationResolver
}
