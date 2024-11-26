package main

import (
	config "go_graphql/config"
	"go_graphql/gql"
	"go_graphql/gql/generated"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize the Gin router
	r := gin.Default()

	// Initialize the database connection
	db := config.InitDB()

	// Migrate the schema (optional, but recommended)
	// database.AutoMigrate(&gqlmodels.Article{}, &gqlmodels.User{}, &gqlmodels.Comment{})

	// Initialize the resolver with the database connection
	resolver := &gql.Resolver{DB: db}

	// Initialize the GraphQL server with the executable schema
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))

	//r.Use(middlewares.AuthMiddleware())
	// Set up routes for GraphQL and Playground
	r.POST("/graphql", gin.WrapH(gqlServer))

	// Start the server on port 8080
	log.Println("Server started at http://localhost:8080/playground")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
