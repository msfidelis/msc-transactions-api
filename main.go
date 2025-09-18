package main

import (
	"log"
	"main/routes/balance"
	"main/routes/cache"
	"main/routes/healthcheck"
	"main/routes/statements"
	"main/routes/transactions"
	"main/routines"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	jsoniter "github.com/json-iterator/go"

	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
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

	app.Get("/statements", statements.GetStatement)
	app.Get("/balance", balance.GetStatement)
	app.Post("/transactions", transactions.NewTransaction)
	app.Get("/transactions/:id_transaction", transactions.DetailTransaction)

	// Cache Examples
	app.Post("/cache/dualwrite/transactions", cache.NewTransactionDualWrite)

	app.Get("/healthcheck", healthcheck.Probe)
	log.Fatal(app.Listen(":8080"))
}
