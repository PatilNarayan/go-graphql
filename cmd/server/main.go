package main

import (
	config "go_graphql/config"
	"go_graphql/gql"
	"go_graphql/gql/generated"
	"go_graphql/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	log := logger.Log

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		logger.AddContext(err).Fatal("Failed to load environment variables")
	}
	log.Info("Environment variables loaded successfully")

	// Initialize Gin router
	r := gin.Default()
	log.Info("Gin router initialized")

	// Initialize database connection
	db := config.InitDB()
	log.Info("Database connection established")

	// Initialize resolver and GraphQL server
	resolver := &gql.Resolver{DB: db}
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	log.Info("GraphQL server initialized")

	// Set up routes
	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))
	r.POST("/graphql", gin.WrapH(gqlServer))
	log.Info("Routes configured")

	// Start server
	log.Info("Starting server at http://localhost:8080/playground")
	if err := r.Run(":8080"); err != nil {
		logger.AddContext(err).Fatal("Server failed to start")
	}
}
