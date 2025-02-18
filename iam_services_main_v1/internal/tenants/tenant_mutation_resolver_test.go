package tenants

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockPermitClient struct {
	mock.Mock
	BaseURL string
}

func TestMain(m *testing.M) {
	//set environment variables
	os.Setenv("PERMIT_PROJECT", "test")
	os.Setenv("PERMIT_ENV", "test")
	os.Setenv("PERMIT_TOKEN", "test")
	os.Setenv("PERMIT_PDP_ENDPOINT", "https://localhost:8080")

	m.Run()
}

// func TestTenantMutationResolver_CreateTenant(t *testing.T) {
// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("failed to connect to database: %v", err)
// 	}

// 	config.DB = db

// 	// Migrate the required schema
// 	err = db.AutoMigrate(&dto.TenantResources{}, &dto.TenantMetadata{}, &dto.Mst_ResourceTypes{}, &dto.TNTRole{}, &dto.TNTPermission{}, &dto.TNTRolePermission{}, &dto.MstRole{}, &dto.MstPermission{}, &dto.MstRolePermission{})
// 	if err != nil {
// 		t.Fatalf("failed to migrate database: %v", err)
// 	}

// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()

// 	// Default responder for unmatched requests
// 	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, `{"error": "unmocked request"}`))

// 	// Register the mock responder for the API endpoint
// 	httpmock.RegisterResponder("POST", "https://localhost:8080/v2/facts/test/test/tenants",
// 		func(req *http.Request) (*http.Response, error) {
// 			resp := httpmock.NewStringResponse(200, `
// 				{
// 					"key": "tenant_123",
// 					"name": "Updated Name",
// 					"status": "success"
// 				}
// 				`)
// 			resp.Header.Add("Content-Type", "application/json")
// 			return resp, nil
// 		},
// 	)

// 	// Seed initial data
// 	mstResType1 := dto.Mst_ResourceTypes{
// 		ResourceTypeID: uuid.New(),
// 		ServiceID:      uuid.New(),
// 		Name:           "Tenant",
// 		RowStatus:      1,
// 		CreatedBy:      uuid.New(),
// 		UpdatedBy:      uuid.New(),
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}
// 	db.Create(&mstResType1)

// 	// Seed initial data
// 	mstResType := dto.Mst_ResourceTypes{
// 		ResourceTypeID: uuid.New(),
// 		ServiceID:      uuid.New(),
// 		Name:           "Role",
// 		RowStatus:      1,
// 		CreatedBy:      uuid.New(),
// 		UpdatedBy:      uuid.New(),
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}
// 	db.Create(&mstResType)

// 	parentOrg := dto.TenantResources{
// 		ResourceID:     uuid.New(),
// 		Name:           "Parent Organization",
// 		ResourceTypeID: mstResType1.ResourceTypeID,
// 		CreatedBy:      uuid.New(),
// 		UpdatedBy:      uuid.New(),
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}
// 	db.Create(&parentOrg)

// 	type args struct {
// 		ctx   context.Context
// 		input models.CreateTenantInput
// 	}
// 	tests := []struct {
// 		name    string
// 		r       *TenantMutationResolver
// 		args    args
// 		want    *models.Tenant
// 		wantErr bool
// 	}{
// 		{
// 			name: "Successful tenant creation",
// 			r:    &TenantMutationResolver{DB: db},
// 			args: args{
// 				ctx: context.Background(),
// 				input: models.CreateTenantInput{
// 					Name: "Test Tenant",
// 					// CreatedBy:  uuid.New(),
// 					Description: ptr.String("A test tenant"),
// 					ContactInfo: &models.ContactInfoInput{
// 						Email:       ptr.String("test@example.com"),
// 						PhoneNumber: ptr.String("1234567890"),
// 						Address: &models.CreateAddressInput{
// 							Street:  ptr.String("123 Test St"),
// 							City:    ptr.String("Test City"),
// 							State:   ptr.String("Test State"),
// 							ZipCode: ptr.String("12345"),
// 							Country: ptr.String("Test Country"),
// 						},
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := tt.r.CreateTenant(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("TenantMutationResolver.CreateTenant() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestTenantMutationResolver_UpdateTenant(t *testing.T) {

// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("failed to connect to database: %v", err)
// 	}

// 	// Migrate the required schema
// 	err = db.AutoMigrate(&dto.TenantResources{}, &dto.TenantMetadata{}, &dto.Mst_ResourceTypes{})
// 	if err != nil {
// 		t.Fatalf("failed to migrate database: %v", err)
// 	}

// 	// Seed initial data
// 	mstResType := dto.Mst_ResourceTypes{
// 		ResourceTypeID: uuid.New(),
// 		ServiceID:      uuid.New(),
// 		Name:           "Tenant",
// 		RowStatus:      1,
// 		CreatedBy:      uuid.New(),
// 		UpdatedBy:      uuid.New(),
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}
// 	db.Create(&mstResType)

// 	existingTenant := dto.TenantResources{
// 		ResourceID:     uuid.New(),
// 		Name:           "Existing Tenant",
// 		ResourceTypeID: mstResType.ResourceTypeID,
// 		CreatedBy:      uuid.New(),
// 		UpdatedBy:      uuid.New(),
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}
// 	db.Create(&existingTenant)

// 	type args struct {
// 		ctx   context.Context
// 		input models.UpdateTenantInput
// 	}
// 	tests := []struct {
// 		name    string
// 		r       *TenantMutationResolver
// 		args    args
// 		want    *models.Tenant
// 		wantErr bool
// 	}{
// 		{
// 			name: "Tenant not found",
// 			r:    &TenantMutationResolver{DB: db},
// 			args: args{
// 				ctx: context.Background(),
// 				input: models.UpdateTenantInput{
// 					ID:          uuid.New(),
// 					Name:        ptr.String("Nonexistent Tenant"),
// 					Description: ptr.String("Should fail to update."),
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Invalid update - empty ID",
// 			r:    &TenantMutationResolver{DB: db},
// 			args: args{
// 				ctx: context.Background(),
// 				input: models.UpdateTenantInput{
// 					ID:          uuid.Nil,
// 					Name:        ptr.String("Invalid Tenant"),
// 					Description: ptr.String("Invalid update due to empty ID."),
// 				},
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			httpmock.Activate()
// 			defer httpmock.DeactivateAndReset()

// 			// Default responder for unmatched requests
// 			httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, `{"error": "unmocked request"}`))

// 			// Register the mock responder for the API endpoint
// 			httpmock.RegisterResponder("PATCH", fmt.Sprintf("https://localhost:8080/v2/facts/test/test/tenants/%s", tt.args.input.ID.String()),
// 				func(req *http.Request) (*http.Response, error) {
// 					resp := httpmock.NewStringResponse(200, ``)
// 					resp.Header.Add("Content-Type", "application/json")
// 					return resp, nil
// 				},
// 			)
// 			_, err := tt.r.UpdateTenant(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("TenantMutationResolver.UpdateTenant() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestTenantMutationResolver_DeleteTenant(t *testing.T) {

// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("failed to connect to database: %v", err)
// 	}

// 	// Migrate the required schema
// 	err = db.AutoMigrate(&dto.TenantResources{}, &dto.TenantMetadata{})
// 	if err != nil {
// 		t.Fatalf("failed to migrate database: %v", err)
// 	}

// 	// Seed initial data
// 	tenantID := uuid.New()
// 	tenantResource := dto.TenantResources{
// 		ResourceID: tenantID,
// 		Name:       "Test Tenant",
// 		CreatedBy:  uuid.New(),
// 		UpdatedBy:  uuid.New(),
// 		CreatedAt:  time.Now(),
// 		UpdatedAt:  time.Now(),
// 	}
// 	tenantMetadata := dto.TenantMetadata{
// 		ResourceID: tenantID,
// 		RowStatus:  1,
// 		Metadata:   json.RawMessage(`{"contactInfo": {"email": "abc", "address": {"city": "String", "state": "String", "street": "String", "country": "String", "zipCode": "String"}, "phoneNumber": "12345"}, "description": "xyz"}`),
// 		CreatedBy:  uuid.New(),
// 		UpdatedBy:  uuid.New(),
// 		CreatedAt:  time.Now(),
// 		UpdatedAt:  time.Now(),
// 	}

// 	db.Create(&tenantResource)
// 	db.Create(&tenantMetadata)

// 	type args struct {
// 		ctx context.Context
// 		id  uuid.UUID
// 	}
// 	tests := []struct {
// 		name    string
// 		r       *TenantMutationResolver
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "Successful tenant deletion",
// 			r:    &TenantMutationResolver{DB: db},
// 			args: args{
// 				ctx: context.Background(),
// 				id:  tenantID,
// 			},
// 			want:    true,
// 			wantErr: false,
// 		},
// 		{
// 			name: "Tenant not found",
// 			r:    &TenantMutationResolver{DB: db},
// 			args: args{
// 				ctx: context.Background(),
// 				id:  uuid.New(), // Nonexistent tenant ID
// 			},
// 			want:    false,
// 			wantErr: true,
// 		},
// 		{
// 			name: "Error during metadata deletion",
// 			r: &TenantMutationResolver{DB: db.Session(&gorm.Session{
// 				SkipDefaultTransaction: true, // Simulate a failure
// 			})},
// 			args: args{
// 				ctx: context.Background(),
// 				id:  uuid.New(),
// 			},
// 			want:    false,
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			httpmock.Activate()
// 			defer httpmock.DeactivateAndReset()

// 			// Default responder for unmatched requests
// 			httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, `{"error": "unmocked request"}`))

// 			// Register the mock responder for the API endpoint
// 			httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://localhost:8080/v2/facts/test/test/tenants/%s", tt.args.id.String()),
// 				func(req *http.Request) (*http.Response, error) {
// 					resp := httpmock.NewStringResponse(200, ``)
// 					resp.Header.Add("Content-Type", "application/json")
// 					return resp, nil
// 				},
// 			)
// 			_, err := tt.r.DeleteTenant(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("TenantMutationResolver.DeleteTenant() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			// if got != tt.want {
// 			// 	t.Errorf("TenantMutationResolver.DeleteTenant() = %v, want %v", got, tt.want)
// 			// }
// 		})
// 	}
// }
