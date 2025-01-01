package tenants

import (
	"github.com/stretchr/testify/mock"
)

type MockPermitClient struct {
	mock.Mock
	BaseURL string
}

// func setupMockServer() *httptest.Server {
// 	// Create a mock server that responds to /v2/facts/test/test/tenants
// 	handler := http.NewServeMux()
// 	handler.HandleFunc("/v2/facts/test/test/tenants", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodPost {
// 			// Respond with a mock JSON response
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusOK)
// 			fmt.Fprintln(w, `{"key": "tenant_123"}`)
// 		} else {
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 		}
// 	})

// 	// Return the mock server
// 	return httptest.NewServer(handler)
// }

// func TestMain(m *testing.M) {
// 	os.Setenv("PERMIT_PROJECT", "test")
// 	os.Setenv("PERMIT_ENV", "test")
// 	os.Setenv("PERMIT_TOKEN", "test")
// 	os.Setenv("PERMIT_PDP_ENDPOINT", "http://localhost:8080")
// 	code := m.Run()

// 	os.Exit(code)
// }
// func (m *MockPermitClient) APIExecute(ctx context.Context, method, endpoint string, body interface{}) (map[string]interface{}, error) {
// 	fullURL := m.BaseURL + endpoint
// 	log.Printf("Mock API call: %s %s", method, fullURL)

// 	args := m.Called(ctx, method, fullURL, body)
// 	if result := args.Get(0); result != nil {
// 		return result.(map[string]interface{}), args.Error(1)
// 	}
// 	return nil, args.Error(1)
// }

// func setupTestDB() *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("failed to connect database: %v", err)
// 	}
// 	db.AutoMigrate(&dto.Tenant{})
// 	return db
// }
// func TestCreateTenant(t *testing.T) {
// 	// Activate httpmock
// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()

// 	// Register the mock responder for the API endpoint
// 	httpmock.RegisterResponder("POST", "http://localhost:8080/v2/facts/test/test/tenants",
// 		func(req *http.Request) (*http.Response, error) {
// 			resp := httpmock.NewStringResponse(200, `
// 					{
// 						"key": "tenant_123",
// 						"status": "success"
// 					}
// 					`)
// 			resp.Header.Add("Content-Type", "application/json")
// 			return resp, nil
// 		},
// 	)

// 	// Setup the test database and the mock permit client
// 	db := setupTestDB()
// 	ctx := context.Background()
// 	// mockPermit := new(MockPermitClient)
// 	resolver := &TenantMutationResolver{DB: db}

// 	// Test input
// 	input := models.TenantInput{
// 		Name:           "Test Tenant",
// 		Description:    ptr.String("Test Tenant Description"),
// 		ParentOrgID:    "org_123",
// 		Metadata:       ptr.String(`{"key": "value"}`),
// 		ParentTenantID: ptr.String("parent_tenant_123"),
// 		CreatedBy:      ptr.String("user_123"),
// 		UpdatedBy:      ptr.String("user_123"),
// 	}

// 	// Call the resolver to create the tenant
// 	tenant, err := resolver.CreateTenant(ctx, input)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, tenant)
// 	assert.Equal(t, "Test Tenant", tenant.Name)
// 	assert.Equal(t, "org_123", tenant.ParentOrgID)
// }

// func TestCreateTenant_MissingName(t *testing.T) {
// 	db := setupTestDB()
// 	ctx := context.Background()

// 	resolver := &TenantMutationResolver{DB: db}

// 	input := models.TenantInput{
// 		Name: "",
// 	}

// 	tenant, err := resolver.CreateTenant(ctx, input)
// 	assert.Error(t, err)
// 	assert.Nil(t, tenant)
// 	assert.Equal(t, "name is required", err.Error())
// }

// func TestUpdateTenant(t *testing.T) {
// 	db := setupTestDB()
// 	ctx := context.Background()
// 	resolver := &TenantMutationResolver{DB: db}

// 	// Seed the database with a tenant
// 	tenant := &dto.Tenant{Name: "Old Name", ParentOrgID: "org_123"}
// 	db.Create(tenant)

// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()

// 	// Default responder for unmatched requests
// 	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, `{"error": "unmocked request"}`))

// 	// Register the mock responder for the API endpoint
// 	httpmock.RegisterResponder("PATCH", fmt.Sprintf("http://localhost:8080/v2/facts/test/test/tenants/%s", tenant.ID),
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
// 	// Test scenarios
// 	t.Run("Valid Input", func(t *testing.T) {
// 		input := models.TenantInput{
// 			Name:        "Updated Name",
// 			ParentOrgID: "org_456",
// 			Description: ptrString("Updated Description"),
// 			Metadata:    ptrString(`{"key": "value"}`),
// 			UpdatedBy:   ptrString("user_123"),
// 		}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, tenant.ID, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, updatedTenant)
// 		assert.Equal(t, "Updated Name", updatedTenant.Name)
// 		assert.Equal(t, "org_456", updatedTenant.ParentOrgID)
// 	})

// 	t.Run("Invalid Tenant ID", func(t *testing.T) {
// 		input := models.TenantInput{
// 			Name: "Updated Name",
// 		}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, "invalid_id", input)
// 		assert.Error(t, err)
// 		assert.Nil(t, updatedTenant)
// 	})

// 	t.Run("HTTPMock Failure", func(t *testing.T) {
// 		// Activate httpmock
// 		httpmock.Activate()
// 		defer httpmock.DeactivateAndReset()

// 		// Register the mock responder
// 		httpmock.RegisterResponder("PATCH", fmt.Sprintf("http://localhost:8080/v2/facts/test/test/tenants/%s", tenant.ID),
// 			httpmock.NewErrorResponder(fmt.Errorf("mock error")),
// 		)

// 		input := models.TenantInput{
// 			Name:        "Updated Name",
// 			ParentOrgID: "org_456",
// 			UpdatedBy:   ptrString("user_123"),
// 		}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, tenant.ID, input)
// 		assert.Error(t, err)
// 		assert.Nil(t, updatedTenant)
// 	})

// 	t.Run("Empty Fields", func(t *testing.T) {
// 		input := models.TenantInput{
// 			Name: "",
// 		}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, tenant.ID, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, updatedTenant)
// 		assert.Equal(t, "Old Name", updatedTenant.Name)
// 	})

// 	t.Run("Invalid Metadata", func(t *testing.T) {
// 		input := models.TenantInput{
// 			Metadata: ptrString("invalid_json"),
// 		}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, tenant.ID, input)
// 		assert.Error(t, err)
// 		assert.Nil(t, updatedTenant)
// 	})

// 	t.Run("Empty Metadata", func(t *testing.T) {
// 		input := models.TenantInput{
// 			Metadata:  ptrString(""),
// 			UpdatedBy: ptrString("user_123"),
// 		}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, tenant.ID, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, updatedTenant)
// 		assert.Equal(t, "", updatedTenant.Metadata)
// 	})

// 	t.Run("Database Unavailable", func(t *testing.T) {
// 		resolver := &TenantMutationResolver{DB: nil}

// 		input := models.TenantInput{
// 			Name:      "Updated Name",
// 			UpdatedBy: ptrString("user_123"),
// 		}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, "some_id", input)
// 		assert.Error(t, err)
// 		assert.Nil(t, updatedTenant)
// 	})

// 	t.Run("Audit Fields", func(t *testing.T) {
// 		input := models.TenantInput{
// 			Name:      "Updated Name",
// 			UpdatedBy: ptrString("user_456"),
// 		}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, tenant.ID, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, updatedTenant)
// 		assert.Equal(t, "Updated Name", updatedTenant.Name)
// 		assert.Equal(t, "user_456", updatedTenant.UpdatedBy)
// 	})

// 	t.Run("No Updates", func(t *testing.T) {
// 		input := models.TenantInput{}

// 		updatedTenant, err := resolver.UpdateTenant(ctx, tenant.ID, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, updatedTenant)
// 		assert.Equal(t, "Old Name", updatedTenant.Name)
// 		assert.Equal(t, "org_123", updatedTenant.ParentOrgID)
// 	})
// }

// func TestDeleteTenant(t *testing.T) {
// 	db := setupTestDB()
// 	ctx := context.Background()

// 	resolver := &TenantMutationResolver{DB: db}

// 	// Seed the database with a tenant
// 	tenant := &dto.Tenant{Name: "Tenant to Delete", RowStatus: 1}
// 	db.Create(tenant)

// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()

// 	// Register the mock responder for the API endpoint
// 	httpmock.RegisterResponder("DELETE", fmt.Sprintf("http://localhost:8080/v2/facts/test/test/tenants/%s", tenant.ID),
// 		func(req *http.Request) (*http.Response, error) {
// 			resp := httpmock.NewStringResponse(200, `
// 				{
// 					"key": "tenant_123",
// 					"status": "success"
// 				}
// 				`)
// 			resp.Header.Add("Content-Type", "application/json")
// 			return resp, nil
// 		},
// 	)

// 	t.Run("Valid Deletion", func(t *testing.T) {
// 		deleted, err := resolver.DeleteTenant(ctx, tenant.ID)
// 		assert.NoError(t, err)
// 		assert.True(t, deleted)

// 		// Ensure the tenant is deleted from the database
// 		var deletedTenant dto.Tenant
// 		result := db.First(&deletedTenant, "id = ?", tenant.ID)
// 		assert.Error(t, result.Error)
// 	})

// 	t.Run("Invalid Tenant ID", func(t *testing.T) {
// 		deleted, err := resolver.DeleteTenant(ctx, "invalid_id")
// 		assert.Error(t, err)
// 		assert.False(t, deleted)
// 	})

// 	t.Run("HTTPMock Failure", func(t *testing.T) {
// 		httpmock.RegisterResponder("DELETE", fmt.Sprintf("http://localhost:8080/v2/facts/test/test/tenants/%s", tenant.ID),
// 			httpmock.NewErrorResponder(fmt.Errorf("mock error")),
// 		)

// 		deleted, err := resolver.DeleteTenant(ctx, tenant.ID)
// 		assert.Error(t, err)
// 		assert.False(t, deleted)
// 	})

// }

// func TestTenant(t *testing.T) {
// 	// Activate httpmock
// 	httpmock.Activate()
// 	defer httpmock.DeactivateAndReset()

// 	// Test Case 1: Valid Input
// 	t.Run("Valid Input", func(t *testing.T) {
// 		// Register mock responder for successful API call
// 		httpmock.RegisterResponder("POST", "http://localhost:8080/v2/facts/test/test/tenants",
// 			func(req *http.Request) (*http.Response, error) {
// 				resp := httpmock.NewStringResponse(200, `{"key": "tenant_123", "status": "success"}`)
// 				resp.Header.Add("Content-Type", "application/json")
// 				return resp, nil
// 			},
// 		)

// 		// Setup the test database and resolver
// 		db := setupTestDB()
// 		ctx := context.Background()
// 		resolver := &TenantMutationResolver{DB: db}

// 		// Test input
// 		input := models.TenantInput{
// 			Name:           "Test Tenant",
// 			Description:    ptr.String("Test Tenant Description"),
// 			ParentOrgID:    "org_123",
// 			Metadata:       ptr.String(`{"key": "value"}`),
// 			ParentTenantID: ptr.String("parent_tenant_123"),
// 			CreatedBy:      ptr.String("user_123"),
// 		}

// 		// Call the resolver to create the tenant
// 		tenant, err := resolver.CreateTenant(ctx, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, tenant)
// 		assert.Equal(t, "Test Tenant", tenant.Name)
// 		assert.Equal(t, "org_123", tenant.ParentOrgID)
// 	})

// 	// Test Case 2: Missing Name
// 	t.Run("Missing Name", func(t *testing.T) {
// 		// Setup the test database and resolver
// 		db := setupTestDB()
// 		ctx := context.Background()
// 		resolver := &TenantMutationResolver{DB: db}

// 		// Test input with missing name
// 		input := models.TenantInput{
// 			Name: "", // Name is missing
// 		}

// 		// Call the resolver to create the tenant
// 		tenant, err := resolver.CreateTenant(ctx, input)
// 		assert.Error(t, err)
// 		assert.Nil(t, tenant)
// 		assert.Equal(t, "name is required", err.Error())
// 	})

// 	// Test Case 3: Permit API Failure
// 	t.Run("Permit API Failure", func(t *testing.T) {
// 		// Register mock responder for failed API call
// 		httpmock.RegisterResponder("POST", "http://localhost:8080/v2/facts/test/test/tenants",
// 			func(req *http.Request) (*http.Response, error) {
// 				resp := httpmock.NewStringResponse(500, `{"status": "error", "message": "internal server error"}`)
// 				resp.Header.Add("Content-Type", "application/json")
// 				return resp, nil
// 			},
// 		)

// 		// Setup the test database and resolver
// 		db := setupTestDB()
// 		ctx := context.Background()
// 		resolver := &TenantMutationResolver{DB: db}

// 		// Test input with valid fields
// 		input := models.TenantInput{
// 			Name: "Test Tenant",
// 		}

// 		// Call the resolver to create the tenant
// 		tenant, err := resolver.CreateTenant(ctx, input)
// 		assert.Error(t, err)
// 		assert.Nil(t, tenant)
// 		fmt.Println("error", err)
// 		assert.Contains(t, err.Error(), "HTTP error: 500")
// 	})

// 	// Test Case 4: No ContactInfoID
// 	t.Run("No ContactInfoID", func(t *testing.T) {
// 		// Register mock responder for successful API call
// 		httpmock.RegisterResponder("POST", "http://localhost:8080/v2/facts/test/test/tenants",
// 			func(req *http.Request) (*http.Response, error) {
// 				resp := httpmock.NewStringResponse(200, `{"key": "tenant_123", "status": "success"}`)
// 				resp.Header.Add("Content-Type", "application/json")
// 				return resp, nil
// 			},
// 		)

// 		// Setup the test database and resolver
// 		db := setupTestDB()
// 		ctx := context.Background()
// 		resolver := &TenantMutationResolver{DB: db}

// 		// Test input without ContactInfoID
// 		input := models.TenantInput{
// 			Name:        "Test Tenant",
// 			ParentOrgID: "org_123",
// 		}

// 		// Call the resolver to create the tenant
// 		tenant, err := resolver.CreateTenant(ctx, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, tenant)
// 	})

// 	// Test Case 10: Empty ParentTenantID
// 	t.Run("Empty ParentTenantID", func(t *testing.T) {
// 		// Register mock responder for successful API call
// 		httpmock.RegisterResponder("POST", "http://localhost:8080/v2/facts/test/test/tenants",
// 			func(req *http.Request) (*http.Response, error) {
// 				resp := httpmock.NewStringResponse(200, `{"key": "tenant_123", "status": "success"}`)
// 				resp.Header.Add("Content-Type", "application/json")
// 				return resp, nil
// 			},
// 		)

// 		// Setup the test database and resolver
// 		db := setupTestDB()
// 		ctx := context.Background()
// 		resolver := &TenantMutationResolver{DB: db}

// 		// Test input with empty ParentTenantID
// 		input := models.TenantInput{
// 			Name:        "Test Tenant",
// 			ParentOrgID: "org_123",
// 		}

// 		// Call the resolver to create the tenant
// 		tenant, err := resolver.CreateTenant(ctx, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, tenant)
// 		assert.Equal(t, "Test Tenant", tenant.Name)
// 		assert.Equal(t, "org_123", tenant.ParentOrgID)
// 	})

// 	// Test Case 11: Optional CreatedBy and UpdatedBy Fields
// 	t.Run("Optional CreatedBy and UpdatedBy Fields", func(t *testing.T) {
// 		// Register mock responder for successful API call
// 		httpmock.RegisterResponder("POST", "http://localhost:8080/v2/facts/test/test/tenants",
// 			func(req *http.Request) (*http.Response, error) {
// 				resp := httpmock.NewStringResponse(200, `{"key": "tenant_123", "status": "success"}`)
// 				resp.Header.Add("Content-Type", "application/json")
// 				return resp, nil
// 			},
// 		)

// 		// Setup the test database and resolver
// 		db := setupTestDB()
// 		ctx := context.Background()
// 		resolver := &TenantMutationResolver{DB: db}

// 		// Test input with missing CreatedBy
// 		input := models.TenantInput{
// 			Name:        "Test Tenant",
// 			ParentOrgID: "org_123",
// 		}

// 		// Call the resolver to create the tenant
// 		tenant, err := resolver.CreateTenant(ctx, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, tenant)
// 	})

// 	// Test Case 12: Tenant Creation with ResourceID
// 	t.Run("Tenant Creation with ResourceID", func(t *testing.T) {
// 		// Register mock responder for successful API call
// 		httpmock.RegisterResponder("POST", "http://localhost:8080/v2/facts/test/test/tenants",
// 			func(req *http.Request) (*http.Response, error) {
// 				resp := httpmock.NewStringResponse(200, `{"key": "tenant_123", "status": "success"}`)
// 				resp.Header.Add("Content-Type", "application/json")
// 				return resp, nil
// 			},
// 		)

// 		// Setup the test database and resolver
// 		db := setupTestDB()
// 		ctx := context.Background()
// 		resolver := &TenantMutationResolver{DB: db}

// 		// Test input with ResourceID
// 		input := models.TenantInput{
// 			Name:        "Test Tenant",
// 			ParentOrgID: "org_123",
// 			ResourceID:  ptr.String("resource_123"),
// 		}

// 		// Call the resolver to create the tenant
// 		tenant, err := resolver.CreateTenant(ctx, input)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, tenant)
// 	})

// }

// func ptrString(s string) *string {
// 	return &s
// }
