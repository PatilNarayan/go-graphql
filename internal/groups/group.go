package groups

import (
	"context"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"

	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

type GroupResolver struct {
	DB *gorm.DB
}

func (r *GroupResolver) Find(ctx context.Context, id string) (*dto.GroupEntity, error) {
	// Implement your logic here
	return nil, nil
}

type GroupQueryResolver struct {
	DB *gorm.DB
}

type GroupMutationResolver struct {
	DB *gorm.DB
}

type GroupFieldResolver struct {
	DB *gorm.DB
}

// CreatedAt implements generated.GroupResolver.
func (g *GroupFieldResolver) CreatedAt(ctx context.Context, obj *dto.GroupEntity) (*string, error) {
	panic("unimplemented")
}

// Tenant implements generated.GroupResolver.
func (g *GroupFieldResolver) Tenant(ctx context.Context, obj *dto.GroupEntity) (*models.Tenant, error) {
	panic("unimplemented")
}

// UpdatedAt implements generated.GroupResolver.
func (g *GroupFieldResolver) UpdatedAt(ctx context.Context, obj *dto.GroupEntity) (*string, error) {
	panic("unimplemented")
}

func (q *GroupQueryResolver) Groups(ctx context.Context) ([]*dto.GroupEntity, error) {
	panic("unimplemented")
}

func (q *GroupQueryResolver) GetGroup(ctx context.Context, id string) (*dto.GroupEntity, error) {
	panic("unimplemented")
}

func (r *GroupMutationResolver) CreateGroup(ctx context.Context, input models.GroupInput) (*dto.GroupEntity, error) {
	panic("unimplemented")
}
func (r *GroupMutationResolver) UpdateGroup(ctx context.Context, id string, input models.GroupInput) (*dto.GroupEntity, error) {
	panic("unimplemented")
}
func (r *GroupMutationResolver) DeleteGroup(ctx context.Context, id string) (bool, error) {
	panic("unimplemented")
}

func (g *GroupFieldResolver) ID(ctx context.Context, obj *dto.GroupEntity) (string, error) {
	panic("unimplemented")
}
