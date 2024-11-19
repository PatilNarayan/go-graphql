package controller

import "go_graphql/permit"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PermitClient *permit.PermitClient
}
