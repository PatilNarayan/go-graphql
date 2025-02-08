package tenants

import (
	"context"
	"encoding/json"
	"go_graphql/gql/models"
	"go_graphql/internal/dto"
	"go_graphql/logger"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestTenantQueryResolver_AllTenants(t *testing.T) {
	// Initialize an in-memory SQLite database for testing
	logger.InitLogger()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Migrate the required schema
	err = db.AutoMigrate(&dto.TenantResource{}, &dto.TenantMetadata{}, &dto.Mst_ResourceTypes{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	mstResType1 := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Tenant",
		RowStatus:      1,
		CreatedBy:      "user1",
		UpdatedBy:      "user1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	db.Create(&mstResType1)

	err = db.Create(&mstResType1).Error

	// Seed test data
	tenant1 := dto.TenantResource{
		ResourceID:     uuid.New(),
		Name:           "Tenant 1",
		ResourceTypeID: mstResType1.ResourceTypeID,
		RowStatus:      1,
		CreatedBy:      "user1",
		UpdatedBy:      "user1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	tenant1Metadata := dto.TenantMetadata{
		ResourceID: tenant1.ResourceID,
		Metadata:   json.RawMessage(`{"contactInfo": {"email": "abc", "address": {"city": "String", "state": "String", "street": "String", "country": "String", "zipCode": "String"}, "phoneNumber": "12345"}, "description": "xyz"}`),
		RowStatus:  1,
		CreatedBy:  "user1",
		UpdatedBy:  "user1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tenant2 := dto.TenantResource{
		ResourceID:     uuid.New(),
		Name:           "Tenant 2",
		ResourceTypeID: uuid.New(),
		RowStatus:      1,
		CreatedBy:      "user1",
		UpdatedBy:      "user1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	tenant2Metadata := dto.TenantMetadata{
		ResourceID: tenant2.ResourceID,
		Metadata:   json.RawMessage(`{"contactInfo": {"email": "abc", "address": {"city": "String", "state": "String", "street": "String", "country": "String", "zipCode": "String"}, "phoneNumber": "12345"}, "description": "xyz"}`),
		RowStatus:  1,
		CreatedBy:  "user1",
		UpdatedBy:  "user1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	db.Create(&tenant1)
	db.Create(&tenant2)
	db.Create(&tenant1Metadata)
	db.Create(&tenant2Metadata)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *TenantQueryResolver
		args    args
		want    []*models.Tenant
		wantErr bool
	}{
		{
			name: "Retrieve all tenants",
			r:    &TenantQueryResolver{DB: db},
			args: args{ctx: context.Background()},
			want: []*models.Tenant{
				{
					ID:   tenant1.ResourceID,
					Name: tenant1.Name,
				},
				{
					ID:   tenant2.ResourceID,
					Name: tenant1.Name,
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty database",
			r:       &TenantQueryResolver{DB: db},
			args:    args{ctx: context.Background()},
			want:    []*models.Tenant{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.AllTenants(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TenantQueryResolver.AllTenants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTenantQueryResolver_GetTenant(t *testing.T) {
	logger.InitLogger()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the required schema
	err = db.AutoMigrate(&dto.TenantResource{}, &dto.TenantMetadata{}, &dto.Mst_ResourceTypes{})
	if err != nil {
		panic(err)
	}

	mstResType1 := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Tenant",
		RowStatus:      1,
		CreatedBy:      "user1",
		UpdatedBy:      "user1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	db.Create(&mstResType1)

	root := dto.TenantResource{
		ResourceID:       uuid.New(),
		Name:             "Tenant 1",
		ResourceTypeID:   mstResType1.ResourceTypeID,
		ParentResourceID: nil,
		RowStatus:        1,
		CreatedBy:        "user1",
		UpdatedBy:        "user1",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	rootMetadata := dto.TenantMetadata{
		ResourceID: root.ResourceID,
		Metadata:   json.RawMessage(`{"contactInfo": {"email": "abc", "address": {"city": "String", "state": "String", "street": "String", "country": "String", "zipCode": "String"}, "phoneNumber": "12345"}, "description": "xyz"}`),
		RowStatus:  1,
		CreatedBy:  "user1",
		UpdatedBy:  "user1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	db.Create(&root)
	db.Create(&rootMetadata)

	// Seed test data
	tenant1 := dto.TenantResource{
		ResourceID:       uuid.New(),
		Name:             "Tenant 1",
		ResourceTypeID:   mstResType1.ResourceTypeID,
		ParentResourceID: &root.ResourceID,
		RowStatus:        1,
		CreatedBy:        "user1",
		UpdatedBy:        "user1",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	tenant1Metadata := dto.TenantMetadata{
		ResourceID: tenant1.ResourceID,
		Metadata:   json.RawMessage(`{"contactInfo": {"email": "abc", "address": {"city": "String", "state": "String", "street": "String", "country": "String", "zipCode": "String"}, "phoneNumber": "12345"}, "description": "xyz"}`),
		RowStatus:  1,
		CreatedBy:  "user1",
		UpdatedBy:  "user1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tenant2 := dto.TenantResource{
		ResourceID:       uuid.New(),
		Name:             "Tenant 2",
		ResourceTypeID:   mstResType1.ResourceTypeID,
		ParentResourceID: &root.ResourceID,
		RowStatus:        1,
		CreatedBy:        "user1",
		UpdatedBy:        "user1",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	tenant2Metadata := dto.TenantMetadata{
		ResourceID: tenant2.ResourceID,
		Metadata:   json.RawMessage(`{"contactInfo": {"email": "abc", "address": {"city": "String", "state": "String", "street": "String", "country": "String", "zipCode": "String"}, "phoneNumber": "12345"}, "description": "xyz"}`),
		RowStatus:  1,
		CreatedBy:  "user1",
		UpdatedBy:  "user1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	db.Create(&tenant1)
	db.Create(&tenant2)
	db.Create(&tenant1Metadata)
	db.Create(&tenant2Metadata)
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		r       *TenantQueryResolver
		args    args
		want    *models.Tenant
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", &TenantQueryResolver{DB: db}, args{ctx: context.Background(), id: tenant1.ResourceID}, &models.Tenant{
			ID:   tenant1.ResourceID,
			Name: tenant1.Name,
		}, false},
		{"test2", &TenantQueryResolver{DB: db}, args{ctx: context.Background(), id: tenant2.ResourceID}, &models.Tenant{
			ID:   tenant2.ResourceID,
			Name: tenant2.Name,
		}, false},
		{"test3", &TenantQueryResolver{DB: db}, args{ctx: context.Background(), id: uuid.Nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.GetTenant(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TenantQueryResolver.GetTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func getDB() *gorm.DB {
	logger.InitLogger()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the required schema
	err = db.AutoMigrate(&dto.TenantResource{}, &dto.TenantMetadata{}, &dto.Mst_ResourceTypes{})
	if err != nil {
		panic(err)
	}

	mstResType1 := dto.Mst_ResourceTypes{
		ResourceTypeID: uuid.New(),
		ServiceID:      uuid.New(),
		Name:           "Tenant",
		RowStatus:      1,
		CreatedBy:      "user1",
		UpdatedBy:      "user1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	db.Create(&mstResType1)

	err = db.Create(&mstResType1).Error

	// Seed test data
	tenant1 := dto.TenantResource{
		ResourceID:     uuid.New(),
		Name:           "Tenant 1",
		ResourceTypeID: mstResType1.ResourceTypeID,
		RowStatus:      1,
		CreatedBy:      "user1",
		UpdatedBy:      "user1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	tenant1Metadata := dto.TenantMetadata{
		ResourceID: tenant1.ResourceID,
		Metadata:   json.RawMessage(`{"contactInfo": {"email": "abc", "address": {"city": "String", "state": "String", "street": "String", "country": "String", "zipCode": "String"}, "phoneNumber": "12345"}, "description": "xyz"}`),
		RowStatus:  1,
		CreatedBy:  "user1",
		UpdatedBy:  "user1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tenant2 := dto.TenantResource{
		ResourceID:     uuid.New(),
		Name:           "Tenant 2",
		ResourceTypeID: uuid.New(),
		RowStatus:      1,
		CreatedBy:      "user1",
		UpdatedBy:      "user1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	tenant2Metadata := dto.TenantMetadata{
		ResourceID: tenant2.ResourceID,
		Metadata:   json.RawMessage(`{"contactInfo": {"email": "abc", "address": {"city": "String", "state": "String", "street": "String", "country": "String", "zipCode": "String"}, "phoneNumber": "12345"}, "description": "xyz"}`),
		RowStatus:  1,
		CreatedBy:  "user1",
		UpdatedBy:  "user1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	db.Create(&tenant1)
	db.Create(&tenant2)
	db.Create(&tenant1Metadata)
	db.Create(&tenant2Metadata)

	return db
}
