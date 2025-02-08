package resource

import (
	"gorm.io/gorm"
)

type ResourceFieldResolver struct {
	DB *gorm.DB
}
