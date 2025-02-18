package bindings

import (
	"context"
	"errors"
	"iam_services_main_v1/gql/models"
	"iam_services_main_v1/internal/dto"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type BindingsQueryResolver struct {
	DB *gorm.DB
}

func (r *BindingsQueryResolver) Binding(ctx context.Context, id uuid.UUID) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className": "binding_query_resolver",
		"method":    "GetBinding",
		"bindingId": id,
	})

	logger.Infof("getbinding request received with Id %v", id)
	if id == uuid.Nil {
		logger.Error("invalid binding id provided")
		return nil, errors.New("binding id is mandatory")
	}
	var binding *dto.TenantRoleAssignments
	if err := r.DB.First(&binding, "resource_id=?", id).Error; err != nil {
		logger.Errorf("error while fetching binding %v", id)
		return nil, err
	}

	principalType := r.FetchPrincipalBasedOnPrincipalId(ctx, binding.PrincipalID)
	createdAt := binding.CreatedAt.String()
	updatedAt := binding.UpdatedAt.String()
	createdBy := binding.CreatedBy
	updatedBy := binding.UpdatedBy
	bindingData := &models.Binding{
		// ID:        binding.ResourceID.String(),
		Name:      binding.Name,
		Version:   binding.Version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
		Role:      &models.Role{ID: binding.RoleID},
	}

	var user *models.User
	if principalType.Name == "User" {
		user = &models.User{
			ID: binding.PrincipalID,
		}
		bindingData.Principal = user
	} else if principalType.Name == "Group" {
		group := &models.Group{
			ID: binding.PrincipalID,
		}
		bindingData.Principal = group
	}

	return nil, nil
}

// AllBindings is the resolver for the allBindings field.
func (r *BindingsQueryResolver) Bindings(ctx context.Context) (models.OperationResult, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className": "binding_query_resolver",
		"method":    "AllBindings",
	})

	logger.Infof("fetch allBindings request received ")
	var roleAssignments []*dto.TenantRoleAssignments
	if err := r.DB.Find(&roleAssignments).Error; err != nil {
		logger.Errorf("error while fetching bindings %v", err)
		return nil, err
	}
	var bindings []*models.Binding

	for _, binding := range roleAssignments {
		// createdAt := binding.CreatedAt.String()
		// updatedAt := binding.UpdatedAt.String()
		// createdBy := binding.CreatedBy
		// updatedBy := binding.UpdatedBy
		principalType := r.FetchPrincipalBasedOnPrincipalId(ctx, binding.PrincipalID)

		bindingData := &models.Binding{
			// ID:      binding.ResourceID.String(),
			Name:    binding.Name,
			Version: binding.Version,
			// CreatedAt: &createdAt,
			// UpdatedAt: &updatedAt,
			// CreatedBy: &createdBy,
			// UpdatedBy: &updatedBy,
			Role: &models.Role{ID: binding.RoleID},
		}
		var user *models.User
		if principalType.Name == "User" {
			user = &models.User{
				ID: binding.PrincipalID,
			}
			bindingData.Principal = user
		} else if principalType.Name == "Group" {
			group := &models.Group{
				ID: binding.PrincipalID,
			}
			bindingData.Principal = group
		}
		bindings = append(bindings, bindingData)
	}
	return nil, nil

}

func (r *BindingsQueryResolver) FetchPrincipalType(ctx context.Context, principalTypeId uuid.UUID) *dto.MstPrincipalTypes {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"class":           "bindings_mutation_resolver",
		"method":          "DeleteBinding",
		"PrincipalTypeId": principalTypeId,
	})
	var principalType *dto.MstPrincipalTypes
	if err := r.DB.First(&principalType, "principal_type_id=?", principalTypeId.String()).Error; err != nil {
		logger.Errorf("permission not present %v", err)
		return nil
	}
	return principalType
}

func (r *BindingsQueryResolver) FetchPrincipalBasedOnPrincipalId(ctx context.Context, principalId uuid.UUID) *dto.MstPrincipalTypes {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"class":           "bindings_query_resolver",
		"method":          "FetchprincipalBasedOnPrincipalId",
		"PrincipalTypeId": principalId,
	})
	var principal *dto.TenantPrincipals
	if err := r.DB.First(&principal, "resource_id=?", principalId).Error; err != nil {
		logger.Errorf("Principal not present %v", err)
		return nil
	}
	return r.FetchPrincipalType(ctx, principal.PrincipalTypeID)
}
