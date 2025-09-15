package main

import (
	clients "main/routes/clients"
	"main/routes/healthcheck"
	"main/routines"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	jsoniter "github.com/json-iterator/go"

	"github.com/ansrivas/fiberprometheus/v2"
)

func main() {
	// Routines
	routines.DatabaseMigration()
	// routines.ClientesMemoryMapping()

	app := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder: jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		Immutable:   true,
		Prefork:     false,
	})

	// Initialize default config
	app.Use(compress.New())

	// Or extend your config for customization
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	// Prometheus
	prometheus := fiberprometheus.New("transaction-api")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// app.Use(pprof.New())

	app.Get("/statements", clients.Statement)
	app.Post("/transactions", clients.NewTransaction)
	app.Get("/healthcheck", healthcheck.Probe)
	app.Listen(":8080")
}
