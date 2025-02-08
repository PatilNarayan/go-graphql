package groups

// type GroupQueryResolver struct {
// 	DB *gorm.DB
// }

// // Tenants resolver for fetching all Tenants
// func (r *GroupQueryResolver) Groups(ctx context.Context) ([]*dto.GroupEntity, error) {
// 	var Groups []*dto.GroupEntity
// 	if err := r.DB.Find(&Groups).Error; err != nil {
// 		return nil, err
// 	}
// 	return Groups, nil
// }

// // GetTenant resolver for fetching a single Tenant by ID
// func (r *GroupQueryResolver) GetGroup(ctx context.Context, id string) (*dto.GroupEntity, error) {
// 	if id == "" {
// 		return nil, errors.New("id cannot be nil")
// 	}

// 	var Group dto.GroupEntity
// 	if err := r.DB.First(&Group, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &Group, nil
// }
