package resource

import (
	"gorm.io/gorm"
)

type ResourceMutationResolver struct {
	DB *gorm.DB
}
