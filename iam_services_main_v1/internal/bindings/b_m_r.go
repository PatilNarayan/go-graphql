package bindings

import (
	"context"
	"iam_services_main_v1/gql/models"

	"gorm.io/gorm"
)

type BindingsMutationResolver struct {
	DB *gorm.DB
}

func (r *BindingsMutationResolver) CreateBinding(ctx context.Context, input models.CreateBindingInput) (*models.Binding, error) {
	// 	logger := log.WithContext(ctx).WithFields(log.Fields{
	// 		"class":       "bindings_mutation_resolver",
	// 		"method":      "CreateBinding",
	// 		"bindingName": input.Name,
	// 	})
	// 	logger.Info("create binding request received")
	// 	if input.Name == "" {
	// 		return nil, errors.New("name is mandatory")
	// 	}
	// 	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)
	// 	if !ok {
	// 		return nil, fmt.Errorf("gin context not found in the request")
	// 	}
	// 	tenantID, exists := ginCtx.Get("tenantID")
	// 	if !exists {
	// 		return nil, fmt.Errorf("tenant id is missing")
	// 	}

	// 	var resourceType *dto.Mst_ResourceTypes
	// 	if err := r.DB.Where(&dto.Mst_ResourceTypes{Name: "Binding"}).First(&resourceType).Error; err != nil {
	// 		logger.Errorf("error while fetching binding for update %v", err)
	// 		return nil, err
	// 	}

	// 	// check whether principal is present or not
	// 	principal := r.IsPrincipalPresent(ctx, input.PrincipalID)
	// 	if input.PrincipalID == "" || principal == nil {
	// 		return nil, errors.New("invalid principal provided")
	// 	}

	// 	// check whether role is present or not
	// 	if input.RoleID == "" || !r.IsRolePresent(ctx, input.RoleID) {
	// 		return nil, errors.New("invalid role provided")
	// 	}

	// 	resourceId := uuid.New()

	// 	// check if a binding is already present
	// 	if r.IsBindingAlreadyPresent(ctx, input.PrincipalID, input.RoleID) {
	// 		return nil, errors.New("role assignment already present")
	// 	}

	// 	parsedTenantId := uuid.MustParse(tenantID.(string))

	// 	currentDate := time.Now()
	// 	bindingDto := &dto.TenantResources{
	// 		ResourceID:       resourceId,
	// 		ResourceTypeID:   resourceType.ResourceTypeID,
	// 		ParentResourceID: nil,
	// 		TenantID:         &parsedTenantId,
	// 		Name:             input.Name,
	// 		RowStatus:        1,
	// 		CreatedAt:        currentDate,
	// 		UpdatedAt:        currentDate,
	// 	}

	// 	//Create binding in resources table
	// 	if err := r.DB.Create(bindingDto).Error; err != nil {
	// 		logger.Errorf("error while creating Binding record %v", err)
	// 		return nil, err
	// 	}

	// 	principalID := uuid.MustParse(input.PrincipalID)
	// 	roleID := uuid.MustParse(input.RoleID)
	// 	tenantBindingDto := &dto.TenantRoleAssignments{
	// 		ResourceID:  resourceId,
	// 		Name:        input.Name,
	// 		Version:     "V1",
	// 		PrincipalID: principalID,
	// 		RoleID:      roleID,
	// 		RowStatus:   1,
	// 		CreatedBy:   input.CreatedBy,
	// 		UpdatedBy:   input.CreatedBy,
	// 		CreatedAt:   time.Now(),
	// 		UpdatedAt:   time.Now(),
	// 	}

	// 	// Create binding in role assignments table
	// 	if err := r.DB.Create(tenantBindingDto).Error; err != nil {
	// 		logger.Error("error while creating Bindings", err)
	// 		return nil, err
	// 	}

	// 	principalType := r.FetchPrincipalType(ctx, principal.PrincipalTypeID)

	// 	createdAt := tenantBindingDto.CreatedAt.String()
	// 	updatedAt := tenantBindingDto.UpdatedAt.String()
	// 	assignment := &models.Binding{
	// 		ID:        resourceId.String(),
	// 		Name:      input.Name,
	// 		CreatedAt: &createdAt,
	// 		UpdatedAt: &updatedAt,
	// 		Principal: &models.Group{},
	// 		Role:      &models.Role{ID: roleID},
	// 		Version:   tenantBindingDto.Version,
	// 	}
	// 	var user *models.User
	// 	if principalType.Name == "User" {
	// 		user = &models.User{
	// 			ID: uuid.MustParse(input.PrincipalID),
	// 		}
	// 		assignment.Principal = user
	// 	} else if principalType.Name == "Group" {
	// 		group := &models.Group{
	// 			ID: uuid.MustParse(input.PrincipalID),
	// 		}
	// 		assignment.Principal = group
	// 	}

	// 	// Create binding in Permit
	// 	pc := permit.NewPermitClient()
	// 	_, err := pc.APIExecute(ctx, "POST", "role_assignments", map[string]interface{}{
	// 		"role":              input.RoleID,
	// 		"tenant":            parsedTenantId,
	// 		"user":              input.PrincipalID,
	// 		"resource_instance": "Tenant:" + parsedTenantId.String(),
	// 	})

	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to create role assignment in permit")
	// 	}

	// 	logger.Info("Binding created successfully")
	// 	return assignment, nil

	return nil, nil
}

func (r *BindingsMutationResolver) UpdateBinding(ctx context.Context, input models.UpdateBindingInput) (*models.Binding, error) {
	// 	logger := log.WithContext(ctx).WithFields(log.Fields{
	// 		"class":     "bindings_mutation_resolver",
	// 		"method":    "CreateBinding",
	// 		"bindingID": input.ID,
	// 	})
	// 	logger.Info("update binding request received")

	// 	if input.ID == "" {
	// 		return nil, errors.New("binding id is mandatory")
	// 	}
	// 	if input.Name == "" {
	// 		return nil, errors.New("name is mandatory")
	// 	}
	// 	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)

	// 	if !ok {
	// 		return nil, fmt.Errorf("gin context not found in the request")
	// 	}

	// 	tenantID, exists := ginCtx.Get("tenantID")
	// 	if !exists {
	// 		return nil, fmt.Errorf("tenant id is missing")
	// 	}
	// 	parsedTenantId := uuid.MustParse(tenantID.(string))

	// 	// check if binding resource is managed
	// 	var resourceType *dto.Mst_ResourceTypes
	// 	if err := r.DB.Where(&dto.Mst_ResourceTypes{Name: "Binding"}).First(&resourceType).Error; err != nil {
	// 		logger.Errorf("error while fetching organization for update %v", err)
	// 		return nil, err
	// 	}

	// 	// check if the binding is in resources table
	// 	var resource *dto.TenantResources
	// 	parsedResourceId := uuid.MustParse(input.ID)
	// 	if err := r.DB.Where(&dto.TenantResources{ResourceID: parsedResourceId, RowStatus: 1}).First(&resource).Error; err != nil {
	// 		logger.Errorf("error while fetching organization for update %v", err)
	// 		return nil, err
	// 	}

	// 	// check if the binding in role assignments table
	// 	var binding *dto.TenantRoleAssignments
	// 	if err := r.DB.First(&binding, "resource_id=?", input.ID).Error; err != nil {
	// 		logger.Errorf("Error while fetching bindings with id %v", input.ID)
	// 		return nil, err
	// 	}

	// 	pc := permit.NewPermitClient()
	// 	deleteRequestBody := map[string]interface{}{
	// 		"role":              binding.RoleID.String(),
	// 		"tenant":            resource.TenantID.String(),
	// 		"user":              binding.PrincipalID.String(),
	// 		"resource_instance": "Tenant:" + resource.TenantID.String(),
	// 	}

	// 	//check if name is different
	// 	if resource.Name != input.Name {
	// 		resource.Name = input.Name
	// 	}

	// 	// check if tenant is different
	// 	if resource.TenantID.String() != tenantID {
	// 		resource.TenantID = &parsedTenantId
	// 	}

	// 	// check if user is present
	// 	principalType := r.IsPrincipalPresent(ctx, input.PrincipalID)

	// 	if input.PrincipalID != binding.PrincipalID.String() && principalType == nil {
	// 		return nil, errors.New("invalid principal provided")
	// 	}

	// 	if !r.IsRolePresent(ctx, input.RoleID) {
	// 		return nil, errors.New("invalid role provided")
	// 	}

	// 	_, err := pc.APIExecute(ctx, "POST", "role_assignments", map[string]interface{}{
	// 		"user":              input.PrincipalID,
	// 		"role":              input.RoleID,
	// 		"tenant":            parsedTenantId,
	// 		"resource_instance": "Tenant:" + parsedTenantId.String(),
	// 	})

	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to create role assignment in permit")
	// 	}

	// 	_, err = pc.APIExecute(ctx, "DELETE", "role_assignments", deleteRequestBody)

	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to delete role assignment in permit")
	// 	}

	// 	// update binding in resource table
	// 	if err := r.DB.Where(&dto.TenantResources{ResourceID: uuid.MustParse(input.ID)}).UpdateColumns(&resource).Error; err != nil {
	// 		logger.Errorf("error while updating binding resource %v", err)
	// 		return nil, err
	// 	}

	// 	binding.Name = input.Name
	// 	binding.PrincipalID = uuid.MustParse(input.PrincipalID)
	// 	binding.RoleID = uuid.MustParse(input.RoleID)
	// 	binding.UpdatedAt = time.Now()
	// 	binding.UpdatedBy = input.UpdatedBy

	// 	//data := map[string]interface{}{"Name": input.Name, "PrincipalID": uuid.MustParse(input.PrincipalID), "ResourceID": uuid.MustParse(input.RoleID), "UpdatedBy": input.UpdatedBy, "UpdatedAt": time.Now()}

	// 	resourceId := uuid.MustParse(input.ID)
	// 	if err := r.DB.Where(&dto.TenantRoleAssignments{ResourceID: resourceId}).UpdateColumns(&binding).Error; err != nil {
	// 		logger.Errorf("error while updating binding %v", err)
	// 		return nil, err
	// 	}

	// 	createdAt := binding.CreatedAt.String()
	// 	updatedAtStr := binding.UpdatedAt.String()
	// 	assignment := &models.Binding{
	// 		ID:        input.ID,
	// 		Name:      input.Name,
	// 		CreatedAt: &createdAt,
	// 		UpdatedAt: &updatedAtStr,
	// 		Role:      &models.Role{ID: uuid.MustParse(input.RoleID)},
	// 		Version:   binding.Version,
	// 	}

	// 	var user *models.User
	// 	if principalType.Name == "User" {
	// 		user = &models.User{
	// 			ID: uuid.MustParse(input.ScopeRefID),
	// 		}
	// 		assignment.Principal = user
	// 	} else if principalType.Name == "Group" {
	// 		group := &models.Group{
	// 			ID: uuid.MustParse(input.ScopeRefID),
	// 		}
	// 		assignment.Principal = group
	// 	}

	// 	return assignment, nil

	return nil, nil
}

// // DeleteBinding is the resolver for the deleteBinding field.
func (r *BindingsMutationResolver) DeleteBinding(ctx context.Context, id string) (bool, error) {
	// 	logger := log.WithContext(ctx).WithFields(log.Fields{
	// 		"class":     "bindings_mutation_resolver",
	// 		"method":    "DeleteBinding",
	// 		"bindingId": id,
	// 	})

	// 	logger.Info("delete binding request received")
	// 	if id == "" {
	// 		logger.Error("invalid id provided. Please provide valid binding id")
	// 		return false, errors.New("id is mandatory")
	// 	}

	// 	ginCtx, ok := ctx.Value("GinContextKey").(*gin.Context)

	// 	if !ok {
	// 		return false, fmt.Errorf("gin context not found in the request")
	// 	}

	// 	tenantID, exists := ginCtx.Get("tenantID")
	// 	if !exists {
	// 		return false, fmt.Errorf("tenant id is missing")
	// 	}
	// 	parsedTenantId := uuid.MustParse(tenantID.(string))

	// 	parsedInputId := uuid.MustParse(id)
	// 	var roleAssignment dto.TenantRoleAssignments
	// 	var tenantResource dto.TenantResources

	// 	if err := r.DB.Where(&dto.TenantRoleAssignments{ResourceID: parsedInputId}).First(&roleAssignment).Error; err != nil {
	// 		logger.Errorf("error while fetching organization for update %v", err)
	// 		return false, err
	// 	}

	// 	if err := r.DB.Where(&dto.TenantResources{ResourceID: parsedInputId}).First(&tenantResource).Error; err != nil {
	// 		logger.Errorf("error while fetching organization for update %v", err)
	// 		return false, err
	// 	}
	// 	pc := permit.NewPermitClient()

	// 	deleteRequestBody := map[string]interface{}{
	// 		"role":              roleAssignment.RoleID.String(),
	// 		"tenant":            parsedTenantId.String(),
	// 		"user":              roleAssignment.PrincipalID.String(),
	// 		"resource_instance": "Tenant:" + parsedTenantId.String(),
	// 	}
	// 	_, err := pc.APIExecute(ctx, "DELETE", "role_assignments", deleteRequestBody)

	// 	if err != nil {
	// 		return false, fmt.Errorf("unable to delete role assignment in permit")
	// 	}

	// 	if err := r.DB.Model(dto.TenantRoleAssignments{}).Where(&dto.TenantRoleAssignments{ResourceID: parsedInputId}).Updates(map[string]interface{}{"RowStatus": 0, "UpdatedBy": "", "UpdatedAt": time.Now()}).Error; err != nil {
	// 		logger.Errorf("error while fetching organization for update %v", err)
	// 		return false, err
	// 	}

	// 	if err := r.DB.Model(dto.TenantResources{}).Where(&dto.TenantResources{ResourceID: parsedInputId}).Updates(map[string]interface{}{"RowStatus": 0, "UpdatedBy": "", "UpdatedAt": time.Now()}).Error; err != nil {
	// 		logger.Errorf("error while fetching organization for update %v", err)
	// 		return false, err
	// 	}

	// 	logger.Info("binding deleted successfully")
	// 	return true, nil
	// }

	// func (r *BindingsMutationResolver) IsPrincipalPresent(ctx context.Context, principalId string) *dto.TenantPrincipals {
	// 	logger := log.WithContext(ctx).WithFields(log.Fields{
	// 		"class":       "bindings_mutation_resolver",
	// 		"method":      "isPrincipalPresent",
	// 		"PrincipalId": principalId,
	// 	})
	// 	var principal *dto.TenantPrincipals
	// 	if err := r.DB.First(&principal, "resource_id=?", principalId).Error; err != nil {
	// 		logger.Errorf("Principal not present %v", err)
	// 		return nil
	// 	}
	// 	return principal
	// }

	// func (r *BindingsMutationResolver) IsRolePresent(ctx context.Context, roleId string) bool {
	// 	logger := log.WithContext(ctx).WithFields(log.Fields{
	// 		"class":  "bindings_mutation_resolver",
	// 		"method": "IsRolePresent",
	// 		"RoleId": roleId,
	// 	})
	// 	var role *dto.TenantRoles
	// 	if err := r.DB.First(&role, "resource_id=?", roleId).Error; err != nil {
	// 		logger.Errorf("invalid role provided %v", err)
	// 		return false
	// 	}
	return true, nil
}

// func (r *BindingsMutationResolver) IsPermissionPresent(ctx context.Context, permissionId string) bool {
// 	logger := log.WithContext(ctx).WithFields(log.Fields{
// 		"class":        "bindings_mutation_resolver",
// 		"method":       "DeleteBinding",
// 		"PermissionId": permissionId,
// 	})
// 	var permissions *dto.TenantPermissions
// 	if err := r.DB.First(&permissions, "permission_id=?", permissionId).Error; err != nil {
// 		logger.Errorf("permission not present %v", err)
// 		return false
// 	}
// 	return true
// }

// func (r *BindingsMutationResolver) IsBindingAlreadyPresent(ctx context.Context, principalId, roleId string) bool {
// 	logger := log.WithContext(ctx).WithFields(log.Fields{
// 		"class":        "bindings_mutation_resolver",
// 		"method":       "DeleteBinding",
// 		"PermissionId": principalId,
// 	})
// 	var binding *dto.TenantRoleAssignments
// 	if err := r.DB.First(&binding, "principal_id=? and role_id=? and row_status=1", principalId, roleId).Error; err != nil {
// 		logger.Errorf("permission not present %v", err)
// 		return false
// 	}
// 	return true
// }

// func (r *BindingsMutationResolver) FetchPrincipalType(ctx context.Context, principalTypeId uuid.UUID) *dto.MstPrincipalTypes {
// 	logger := log.WithContext(ctx).WithFields(log.Fields{
// 		"class":           "bindings_mutation_resolver",
// 		"method":          "DeleteBinding",
// 		"PrincipalTypeId": principalTypeId,
// 	})
// 	var principalType *dto.MstPrincipalTypes
// 	if err := r.DB.First(&principalType, "principal_type_id=?", principalTypeId.String()).Error; err != nil {
// 		logger.Errorf("permission not present %v", err)
// 		return nil
// 	}
// 	return principalType
// }
