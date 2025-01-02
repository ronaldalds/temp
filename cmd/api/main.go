package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ronaldalds/base-go-api/internal/database"
	"github.com/ronaldalds/base-go-api/internal/routes"
	"github.com/ronaldalds/base-go-api/internal/server"
	"github.com/ronaldalds/base-go-api/internal/settings"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *fiber.App, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	settings.Load()
	database.DbLoad()
	app := server.New()
	routes := routes.NewRouter(app)
	routes.RegisterFiberRoutes()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		err := app.Listen(fmt.Sprintf(":%d", settings.Env.Port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(app, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
