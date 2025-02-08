package main

import (
	config "go_graphql/config"
	gql "go_graphql/gql"
	"go_graphql/gql/generated"
	"go_graphql/internal/middlewares"
	"go_graphql/pkg/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	// Initialize Gin router
	r := gin.Default()

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		logger.LogFatal("Failed to load environment variables")
	}

	// Initialize database connection
	db := config.InitDB()

	// Initialize resolver and GraphQL server
	resolver := &gql.Resolver{DB: db}
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	// Set up routes

	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))

	r.Use(middlewares.AuthMiddleware())
	r.Use(middlewares.GinContextToContextMiddleware())
	r.Use(middlewares.RequestLogger())

	r.POST("/graphql", func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		gin.WrapH(gqlServer)(ctx)
	})

	// Start server
	if err := r.Run(":8080"); err != nil {
		//logger.AddContext(err).Fatal("Server failed to start")
	}
}
