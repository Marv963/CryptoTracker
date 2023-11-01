package jobs

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/Marv963/CryptoTracker/app/pkg/model/symbol"
)

// calculateArbitrage is a method that operates directly on the current data instead of fetching it from Redis.
func (c *JobController) calculateArbitrage(ctx context.Context, pipe redis.Pipeliner, symbolData *symbol.Symbol) {
	// Stelle sicher, dass sowohl priceBitstamp als auch priceKraken nicht nil sind
	if symbolData.PriceBitstamp == nil || symbolData.PriceKraken == nil {
		return
	}
	// Determine the lower and higher price
	if *symbolData.PriceBitstamp < *symbolData.PriceKraken {
		symbolData.LowestPrice = symbolData.PriceBitstamp
		symbolData.HighestPrice = symbolData.PriceKraken
		symbolData.LowestPriceExchange = "bitstamp"
		symbolData.HighestPriceExchange = "kraken"
	} else {
		symbolData.LowestPrice = symbolData.PriceKraken
		symbolData.HighestPrice = symbolData.PriceBitstamp
		symbolData.LowestPriceExchange = "kraken"
		symbolData.HighestPriceExchange = "bitstamp"
	}

	// Calculate the percent difference
	arbitrage := (*symbolData.HighestPrice - *symbolData.LowestPrice) / *symbolData.LowestPrice * 100
	symbolData.Arbitrage = &arbitrage
}
