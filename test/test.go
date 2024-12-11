package main

import (
	"context"
	"go_graphql/permit"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize the client
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	pc := permit.NewPermitClient()

	tenantData := map[string]interface{}{
		"key":  "example-tenant1",
		"name": "Example Tenant1",
	}

	tenant, err := pc.APIExecute(context.Background(), "POST", "tenants", tenantData)
	if err != nil {
		log.Fatalf("Error creating tenant: %v", err)
	}

	log.Printf("Tenant created: %v", tenant)
}
