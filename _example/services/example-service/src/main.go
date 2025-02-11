package main

import (
	"context"
	"fmt"

	"github.com/dennis-dko/go-toolkit/database"
	"github.com/dennis-dko/go-toolkit/example/src/config"
	"github.com/dennis-dko/go-toolkit/example/src/router"

	"github.com/dennis-dko/go-toolkit/example/docs"

	"github.com/dennis-dko/go-toolkit/server"
)

//	@title			EXAMPLE
//	@version		1.0
//	@description	Example Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam.
//	@contact.name	Example
//	@contact.url	https://www.example.com/
//	@contact.email	example@example.de
//
// @BasePath	/
func main() {
	// Create background context for all context needs
	ctx := context.Background()

	// Initialize configuration
	cfg, envFiles := config.Init(ctx)

	// Initialize MongoDB
	mongoDb, mgDisconnect := database.MongoDBInit(ctx, &cfg.Persistence.MongoDB)
	defer mgDisconnect()

	// Initialize Postgres
	postgres, psDisconnect := database.PostgresInit(ctx, &cfg.Persistence.Postgres)
	defer psDisconnect()

	// Set swagger info host
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// Initialize server
	app := server.New(ctx, &cfg.Server, envFiles)

	// Initialize router
	router.Init(app, cfg, postgres, mongoDb)

	// Start application
	app.Start()
}
