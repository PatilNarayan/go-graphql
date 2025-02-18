package tenants

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of *gorm.DB
type MockDB struct {
	mock.Mock
	*gorm.DB // Embed *gorm.DB to make it compatible with *gorm.DB
}

// func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
// 	called := m.Called(query, args)
// 	return called.Get(0).(*gorm.DB)
// }

// func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
// 	called := m.Called(dest)
// 	return called.Get(0).(*gorm.DB)
// }

// func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
// 	called := m.Called(dest)
// 	return called.Get(0).(*gorm.DB)
// }

// func TestGetTenant(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		tenantID       uuid.UUID
// 		setupMock      func(*MockDB)
// 		expectedError  error
// 		expectedTenant *models.Tenant
// 	}{
// 		{
// 			name:     "successful tenant retrieval",
// 			tenantID: uuid.New(),
// 			setupMock: func(db *MockDB) {
// 				resourceType := &dto.Mst_ResourceTypes{
// 					ResourceTypeID: uuid.New(),
// 					Name:           "Tenant",
// 				}

// 				tenantResource := &dto.TenantResources{
// 					ResourceID:     uuid.New(),
// 					ResourceTypeID: resourceType.ResourceTypeID,
// 					Name:           "Test Tenant",
// 				}

// 				metadata := &dto.TenantMetadata{
// 					ResourceID: tenantResource.ResourceID,
// 					Metadata: json.RawMessage(`{
// 						"description": "Test Description",
// 						"contactInfo": {
// 							"email": "test@example.com",
// 							"phone": "123-456-7890"
// 						}
// 					}`),
// 				}

// 				// Mock resource type query
// 				db.On("Where", mock.Anything, mock.Anything).Return(db)
// 				db.On("First", mock.AnythingOfType("*dto.Mst_ResourceTypes")).
// 					Run(func(args mock.Arguments) {
// 						arg := args.Get(0).(*dto.Mst_ResourceTypes)
// 						*arg = *resourceType
// 					}).
// 					Return(&gorm.DB{Error: nil})

// 				// Mock tenant resource query
// 				db.On("Where", mock.Anything, mock.Anything).Return(db)
// 				db.On("First", mock.AnythingOfType("*dto.TenantResources")).
// 					Run(func(args mock.Arguments) {
// 						arg := args.Get(0).(*dto.TenantResources)
// 						*arg = *tenantResource
// 					}).
// 					Return(&gorm.DB{Error: nil})

// 				// Mock metadata query
// 				db.On("Where", mock.Anything, mock.Anything).Return(db)
// 				db.On("First", mock.AnythingOfType("*dto.TenantMetadata")).
// 					Run(func(args mock.Arguments) {
// 						arg := args.Get(0).(*dto.TenantMetadata)
// 						*arg = *metadata
// 					}).
// 					Return(&gorm.DB{Error: nil})
// 			},
// 			expectedError: nil,
// 			expectedTenant: &models.Tenant{
// 				Name:        "Test Tenant",
// 				Description: ptr.String("Test Description"),
// 				ContactInfo: &models.ContactInfo{
// 					Email:       ptr.String("test@example.com"),
// 					PhoneNumber: ptr.String("123-456-7890"),
// 				},
// 			},
// 		},
// 		{
// 			name:     "tenant not found",
// 			tenantID: uuid.New(),
// 			setupMock: func(db *MockDB) {
// 				// Mock resource type query
// 				db.On("Where", mock.Anything, mock.Anything).Return(db)
// 				db.On("First", mock.AnythingOfType("*dto.Mst_ResourceTypes")).
// 					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
// 			},
// 			expectedError:  ErrResourceTypeNotFound,
// 			expectedTenant: nil,
// 		},
// 		{
// 			name:           "nil tenant ID",
// 			tenantID:       uuid.Nil,
// 			setupMock:      func(db *MockDB) {},
// 			expectedError:  ErrTenantIDRequired,
// 			expectedTenant: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := &MockDB{}
// 			tt.setupMock(mockDB)

// 			resolver := &TenantQueryResolver{
// 				DB: mockDB.DB, // Use the mock directly
// 			}

// 			_, err := resolver.GetTenant(context.Background(), tt.tenantID)

// 			if tt.expectedError != nil {
// 				assert.ErrorIs(t, err, tt.expectedError)
// 			} else {
// 				assert.NoError(t, err)
// 				// assert.Equal(t, tt.expectedTenant.Name, tenant.IsOperationResult())
// 				// assert.Equal(t, tt.expectedTenant.Description, tenant.Description)
// 				if tt.expectedTenant.ContactInfo != nil {
// 					// assert.Equal(t, tt.expectedTenant.ContactInfo.Email, tenant.ContactInfo.Email)
// 					// assert.Equal(t, tt.expectedTenant.ContactInfo.PhoneNumber, tenant.ContactInfo.PhoneNumber)
// 				}
// 			}
// 		})
// 	}
// }

// func TestAllTenants(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		setupMock     func(*MockDB)
// 		expectedError error
// 		expectedCount int
// 	}{
// 		{
// 			name: "successful tenants retrieval",
// 			setupMock: func(db *MockDB) {
// 				resourceType := &dto.Mst_ResourceTypes{
// 					ResourceTypeID: uuid.New(),
// 					Name:           "Tenant",
// 				}

// 				tenantResources := []dto.TenantResources{
// 					{
// 						ResourceID:     uuid.New(),
// 						ResourceTypeID: resourceType.ResourceTypeID,
// 						Name:           "Tenant 1",
// 					},
// 					{
// 						ResourceID:     uuid.New(),
// 						ResourceTypeID: resourceType.ResourceTypeID,
// 						Name:           "Tenant 2",
// 					},
// 				}

// 				// Mock resource type query
// 				db.On("Where", mock.Anything, mock.Anything).Return(db)
// 				db.On("First", mock.AnythingOfType("*dto.Mst_ResourceTypes")).
// 					Run(func(args mock.Arguments) {
// 						arg := args.Get(0).(*dto.Mst_ResourceTypes)
// 						*arg = *resourceType
// 					}).
// 					Return(&gorm.DB{Error: nil})

// 				// Mock tenants query
// 				db.On("Where", mock.Anything).Return(db)
// 				db.On("Find", mock.AnythingOfType("*[]dto.TenantResources")).
// 					Run(func(args mock.Arguments) {
// 						arg := args.Get(0).(*[]dto.TenantResources)
// 						*arg = tenantResources
// 					}).
// 					Return(&gorm.DB{Error: nil})

// 				// Mock metadata queries
// 				db.On("Where", mock.Anything, mock.Anything).Return(db)
// 				db.On("First", mock.AnythingOfType("*dto.TenantMetadata")).
// 					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
// 			},
// 			expectedError: nil,
// 			expectedCount: 2,
// 		},
// 		{
// 			name: "resource type not found",
// 			setupMock: func(db *MockDB) {
// 				db.On("Where", mock.Anything, mock.Anything).Return(db)
// 				db.On("First", mock.AnythingOfType("*dto.Mst_ResourceTypes")).
// 					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
// 			},
// 			expectedError: ErrResourceTypeNotFound,
// 			expectedCount: 0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := &MockDB{}
// 			tt.setupMock(mockDB)

// 			resolver := &TenantQueryResolver{
// 				DB: mockDB.DB, // Use the mock directly
// 			}

// 			_, err := resolver.AllTenants(context.Background())

// 			if tt.expectedError != nil {
// 				assert.ErrorIs(t, err, tt.expectedError)
// 			} else {
// 				assert.NoError(t, err)
// 				// assert.Equal(t, tt.expectedCount, len(tenants))
// 			}
// 		})
// 	}
// }
