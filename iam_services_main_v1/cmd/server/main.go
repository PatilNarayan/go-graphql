package main

import (
	config "iam_services_main_v1/config"
	"iam_services_main_v1/gql"
	"iam_services_main_v1/gql/generated"
	"iam_services_main_v1/internal/middlewares"
	"iam_services_main_v1/internal/permit"
	"iam_services_main_v1/pkg/logger"
	"log"

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
		log.Fatal(err)
	}

	// Initialize database connection
	db := config.InitDB()

	//Initialize permit
	pc := permit.NewPermitClient()

	// Initialize resolver and GraphQL server
	resolver := &gql.Resolver{DB: db, PC: pc}
	gqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	// Set custom error formatting globally
	// gqlServer.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
	// 	return utils.FormatError(ctx, err) // Call your custom error formatting function
	// })

	// Set up routes

	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))

	r.Use(middlewares.AuthMiddleware())
	r.Use(middlewares.GinContextToContextMiddleware())
	// r.Use(middlewares.RequestLogger())

	r.POST("/graphql", func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		gin.WrapH(gqlServer)(ctx)
	})

	// Start server
	if err := r.Run(":8080"); err != nil {
		//logger.AddContext(err).Fatal("Server failed to start")
	}
}
