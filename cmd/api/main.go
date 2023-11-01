package main

import (
	"log"
	"os"

	"github.com/Marv963/CryptoTracker/app/cmd/api/routes"
	"github.com/Marv963/CryptoTracker/app/pkg/appcontext"
	"github.com/Marv963/CryptoTracker/app/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg, err := config.Read("")
	if err != nil {
		log.Fatalf("%s: %s", "Failed to read configuration", err)
		os.Exit(2)
	}

	appContext, err := appcontext.NewAppContext(cfg)
	if err != nil {
		appContext.Logger.Fatalf("Failed to initialize app context: %s", err)
		os.Exit(2)
	}
	defer appContext.Close()

	app := fiber.New()

	// Cors
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins: cfg.Frontend.URL,
		AllowMethods: "GET",
	}))

	cfg = nil

	routes.SetupRoutes(app, appContext)

	log.Fatal(app.Listen(":3003"))
}
