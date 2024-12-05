package main

import (
	"context"
	"go_graphql/permit"
	"log"

	"github.com/joho/godotenv"
	"github.com/rs/xid"
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

	tenant, err := pc.CreateTenant(context.Background(), tenantData)
	if err != nil {
		log.Fatalf("Failed to create tenant: %v", err)
	}

	CreateResourceInstance := map[string]interface{}{
		"key":    xid.New().String(),
		"tenant": tenant["key"],
	}

	_, err = pc.CreateResourceInstance(context.Background(), CreateResourceInstance)
	if err != nil {
		log.Fatalf("Failed to create resource instance: %v", err)
	}
	log.Printf("Tenant created successfully: %+v", tenant)

	log.Printf("Tenant created successfully: %+v", pc)
}
