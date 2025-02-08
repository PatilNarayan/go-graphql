package clientorganizationunit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"iam_services_main_v1/gql/models"
	dto "iam_services_main_v1/internal/dto"
	model "iam_services_main_v1/internal/models"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ClientOrganizationUnitMutationResolver struct {
	DB *gorm.DB
}

func (r *ClientOrganizationUnitMutationResolver) CreateClientOrganizationUnit(ctx context.Context, input models.CreateClientOrganizationUnitInput) (*models.ClientOrganizationUnit, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className":  "organization_mutation_resolver",
		"methodName": "CreateClientOrganizationUnit",
	})
	logger.Info("create clientOrganizationUnit")
	if input.Name == "" {
		return nil, errors.New("organization unit name is mandatory")
	}

	if input.TenantID == "" {
		return nil, errors.New("organization must be created under a tenant")
	}

	var resourceType *dto.Mst_ResourceTypes
	if err := r.DB.Where(&dto.Mst_ResourceTypes{Name: "ClientOrganizationUnit"}).First(&resourceType).Error; err != nil {
		logger.Errorf("error while fetching organization for update %v", err)
		return nil, err
	}

	var parentOrg *dto.TenantResource
	parentOrgId := uuid.MustParse(input.ParentOrgID)
	if err := r.DB.Where(&dto.TenantResource{ResourceID: parentOrgId}).First(&parentOrg).Error; err != nil {
		logger.Errorf("error while fetching organization for update %v", err)
		return nil, err
	}

	tenantId := uuid.MustParse(input.TenantID)
	if input.TenantID != "" {

		if err := r.DB.Where(&dto.TenantResource{ResourceID: tenantId}).First(&parentOrg).Error; err != nil {
			logger.Errorf("error while fetching organization for update %v", err)
			return nil, err
		}
	}

	resourceId := uuid.New()
	currentDate := time.Now()
	corgDto := &dto.TenantResource{
		ResourceID:       resourceId,
		ResourceTypeID:   resourceType.ResourceTypeID,
		ParentResourceID: &parentOrg.ResourceID,
		TenantID:         &tenantId,
		Name:             input.Name,
		RowStatus:        1,
		CreatedAt:        currentDate,
		UpdatedAt:        currentDate,
	}

	if err := r.DB.Create(corgDto).Error; err != nil {
		logger.Errorf("error while creating ClientOrganization record %v", err)
		return nil, err
	}

	clientOrg := model.ClientOrganizationUnit{
		ResourceID:       resourceId,
		ResourceTypeID:   resourceType.ResourceTypeID,
		ParentResourceID: parentOrg.ResourceID,
		Name:             input.Name,
		Description:      *input.Description,
		RowStatus:        1,
		CreatedAt:        currentDate,
		UpdatedAt:        currentDate,
	}

	jsonData, err := json.Marshal(clientOrg)
	if err != nil {
		logger.Errorf("error while unmarshalling client organization %v", err)
	}

	metadataDto := &dto.TenantMetadata{
		ResourceID: resourceId,
		Metadata:   jsonData,
		CreatedAt:  currentDate,
		UpdatedAt:  currentDate,
	}

	if err = r.DB.Create(metadataDto).Error; err != nil {
		logger.Errorf("error while creating metadata record %v", err)
		return nil, err
	}

	response := &models.ClientOrganizationUnit{
		ID:          resourceId,
		Name:        corgDto.Name,
		Description: input.Description,
		ParentOrg:   nil,
		Tenant:      nil,
		CreatedAt:   currentDate.String(),
		UpdatedAt:   nil,
	}

	return response, nil
}

func (r *ClientOrganizationUnitMutationResolver) UpdateClientOrganizationUnit(ctx context.Context, input models.UpdateClientOrganizationUnitInput) (*models.ClientOrganizationUnit, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className":  "client_organization_mutation_resolver",
		"methodName": "UpdateClientOrganizationUnit",
	})
	var zeroUUID uuid.UUID
	if input.ID == zeroUUID {
		return nil, errors.New("id is mandatory for update")
	}
	parsedUuid := uuid.MustParse(*input.TenantID)

	var resource *dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{ResourceID: input.ID, ParentResourceID: &parsedUuid, RowStatus: 1}).First(&resource).Error; err != nil {
		logger.Errorf("error while fetching organization for update %v", err)
		return nil, err
	}

	if resource.Name != *input.Name {
		resource.Name = *input.Name
	}

	parentOrgId := uuid.MustParse(*input.ParentOrgID)
	if resource.ParentResourceID != &parentOrgId {
		resource.ParentResourceID = &parentOrgId
	}

	resource.UpdatedAt = time.Now()

	if err := r.DB.Where(&dto.TenantResource{ResourceID: input.ID}).UpdateColumns(&resource).Error; err != nil {
		logger.Errorf("error while updating client organization unit %v", err)
		return nil, err
	}

	var resourceMetadata dto.TenantMetadata
	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: resource.ResourceID}).First(&resourceMetadata).Error; err != nil {
		return nil, errors.New("unable to find resource metadata")
	}

	// Unmarshal the existing metadata
	metadata := map[string]interface{}{}
	if err := json.Unmarshal([]byte(resourceMetadata.Metadata), &metadata); err != nil {
		return nil, errors.New("unable to unmarshal data")
	}

	if input.Name != metadata["name"] {
		metadata["name"] = input.Name
	}

	if input.Description != metadata["description"] {
		metadata["description"] = input.Description
	}

	if input.TenantID != metadata["tenantId"] {
		metadata["tenantId"] = input.TenantID
	}

	if input.ParentOrgID != metadata["parentOrgId"] {
		metadata["parentOrgId"] = input.ParentOrgID
	}

	updatedMetadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, errors.New("failed to marshal updated metadata")
	}
	resourceMetadata.Metadata = updatedMetadataJSON
	resourceMetadata.UpdatedAt = time.Now()

	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: input.ID}).UpdateColumns(&resourceMetadata).Error; err != nil {
		return nil, fmt.Errorf("failed to create tenant metadata: %w", err)
	}

	updatedAt := resource.UpdatedAt.String()
	response := &models.ClientOrganizationUnit{
		ID:          resource.ResourceID,
		Name:        resource.Name,
		Description: input.Description,
		ParentOrg:   nil,
		Tenant:      nil,
		CreatedAt:   resource.CreatedAt.String(),
		UpdatedAt:   &updatedAt,
	}
	return response, nil
}

func (r *ClientOrganizationUnitMutationResolver) DeleteClientOrganizationUnit(ctx context.Context, id uuid.UUID) (bool, error) {
	logger := log.WithContext(ctx).WithFields(log.Fields{
		"className":  "client_organization_mutation_resolver",
		"methodName": "DeleteClientOrganizationUnit",
	})

	var zeroUUID uuid.UUID
	if id == zeroUUID {
		return false, errors.New("id is mandatory for delete")
	}

	var resource *dto.TenantResource
	if err := r.DB.Where(&dto.TenantResource{ResourceID: id}).First(&resource).Error; err != nil {
		logger.Errorf("error while fetching organization for update %v", err)
		return false, err
	}

	resource.RowStatus = 0
	resource.UpdatedAt = time.Now()

	if err := r.DB.Model(dto.TenantResource{}).Where(&dto.TenantResource{ResourceID: id}).Updates(map[string]interface{}{"RowStatus": 0, "UpdatedBy": "", "UpdatedAt": time.Now()}).Error; err != nil {
		logger.Errorf("error while fetching organization for update %v", err)
		return false, err
	}

	if err := r.DB.Where(&dto.TenantMetadata{ResourceID: id}).Delete(&dto.TenantMetadata{}).Error; err != nil {
		return false, errors.New("failed to delete tenant metadata")
	}

	return true, nil
}
