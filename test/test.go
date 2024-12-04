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
		"key":  "example-tenant",
		"name": "Example Tenant",
	}

	tenant, err := pc.CreateTenant(context.Background(), tenantData)
	if err != nil {
		log.Fatalf("Failed to create tenant: %v", err)
	}

	log.Printf("Tenant created successfully: %+v", tenant)

	log.Printf("Tenant created successfully: %+v", pc)
}
