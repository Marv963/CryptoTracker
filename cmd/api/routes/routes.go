package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Marv963/CryptoTracker/app/cmd/api/handler"
	"github.com/Marv963/CryptoTracker/app/pkg/appcontext"
)

func SetupRoutes(app *fiber.App, appCtx *appcontext.AppContext) {
	// Symbols
	symbolController := handler.NewSymbolController(appCtx)
	app.Get("/api/symbols/:page<int>?", symbolController.GetSymbols)
	app.Get("/api/symbolhistory/:pair<str>", symbolController.GetSymbolHistory)
	app.Get("/api/arbitrage", symbolController.GetArbitrage)
}
