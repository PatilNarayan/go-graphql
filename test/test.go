package main

import (
	"log"
	"os"

	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/permit"
)

func NewPermitClient() (*permit.Client, error) {
	// Create config with provided API token
	permitConfig := config.NewConfigBuilder(os.Getenv("PERMIT_TOKEN")).Build()

	// Initialize the Permit client
	client := permit.New(permitConfig)

	return client, nil
}

func main() {
	// Initialize the client
	pc, err := NewPermitClient()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Tenant created successfully: %+v", pc)
}
