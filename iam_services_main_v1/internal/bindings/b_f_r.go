// bindings/b_f_r.go
package bindings

import (
	"context"
	"iam_services_main_v1/gql/models"

	"gorm.io/gorm"
)

type BindingsFieldResolver struct {
	DB *gorm.DB
}

func (b *BindingsFieldResolver) Principal(ctx context.Context, obj *models.Binding) (models.Principal, error) {
	return nil, nil
}
func (b *BindingsFieldResolver) Role(ctx context.Context, obj *models.Binding) (*models.Role, error) {
	return nil, nil
}
func (b *BindingsFieldResolver) ScopeRef(ctx context.Context, obj *models.Binding) (models.Resource, error) {
	return nil, nil
}
