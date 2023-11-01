package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/redis/go-redis/v9"

	"github.com/Marv963/CryptoTracker/app/pkg/model/symbol"
	"github.com/Marv963/CryptoTracker/app/pkg/workerpool"
)

func (c *JobController) fetchAllSymbolHashes() (map[string]symbol.Symbol, error) {
	// Create a context.Context with a timeout of 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// Make sure to clean up the context at the end
	defer cancel()

	// Retrieve all keys that start with "hash:"
	keys, err := c.appContext.Redis.Keys(ctx, "hash:*").Result()
	if err != nil {
		return nil, fmt.Errorf("could not fetch key: %v", err)
	}

	allSymbolData := make(map[string]symbol.Symbol, len(keys))

	for _, key := range keys {
		val, err := c.appContext.Redis.Get(ctx, key).Result()
		if err != nil {
			c.appContext.Logger.Printf("Failed to get value for key %s: %v \n", key, err)
			continue
		}

		var data symbol.Symbol
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			c.appContext.Logger.Printf("Failed to unmarshal json data: %v \n", err)
			continue
		}

		// Assume key is something like "hash:BTC/USD" and you want "BTC/USD" as the key in your map.
		pair := strings.TrimPrefix(key, "hash:")
		allSymbolData[pair] = data
	}

	return allSymbolData, nil
}

func (c *JobController) FetchSymbolPriceChanges() {
	// Logik zum Abrufen der Symbole und Daten von InfluxDB

	// Verwenden die a.appContext.QueryInfluxDB, um die Daten von InfluxDB abzurufen
	result, err := c.appContext.InfluxClient.Query("queries/getSymbolPriceAndChange.flux")
	if err != nil {
		c.appContext.Logger.Printf("Failed to query get_symbols: %s \n", err)
	}

	// Erhalte alle SymbolHashes um preisänderungen zu identifizieren
	symbolHashes, err := c.fetchAllSymbolHashes()
	if err != nil {
		c.appContext.Logger.Printf("Failed to fetch alle Symbol hashes err: %s \n", err)
	}

	// Create a context.Context with a timeout of 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// Make sure to clean up the context at the end
	defer cancel()

	// Alle Symbol Hashes aus Redis abrufen, um Preisänderungen an den websocket zu senden

	// Create a pipeline for redis HSet
	pipe := c.appContext.Redis.Pipeline()

	// Iterate over the QueryTableResult and extract the data
	for result.Next() {
		pair, ok := result.Record().ValueByKey("pair").(string)
		if !ok {
			c.appContext.Logger.Printf("Error: Value pair is not a string \n")
			continue
		}

		priceLast, ok := c.extractFloatValue(result, "priceLast")
		if !ok {
			continue
		}

		// Check if the price has changed.
		oldSymbolData, exists := symbolHashes[pair]
		if exists {
			previousPrice := oldSymbolData.Price
			if priceLast == previousPrice {
				// The price hasn't changed, so it doesn't need to be updated either
				continue
			}
		}

		price1hAgo, ok := c.extractFloatValue(result, "price1hAgo")
		if !ok {
			continue
		}

		price1dAgo, ok := c.extractFloatValue(result, "price1dAgo")
		if !ok {
			continue
		}

		price7dAgo, ok := c.extractFloatValue(result, "price7dAgo")
		if !ok {
			continue
		}

		priceChange1h, ok := c.extractFloatValue(result, "priceChange1h")
		if !ok {
			continue
		}

		priceChange1d, ok := c.extractFloatValue(result, "priceChange1d")
		if !ok {
			continue
		}

		priceChange7d, ok := c.extractFloatValue(result, "priceChange7d")
		if !ok {
			continue
		}

		priceBitstamp, ok := c.extractOptionalFloat(result, "priceBitstamp")
		if !ok {
			continue
		}

		priceKraken, ok := c.extractOptionalFloat(result, "priceKraken")
		if !ok {
			continue
		}

		bidKraken, ok := c.extractOptionalFloat(result, "bidKraken")
		if !ok {
			continue
		}

		askKraken, ok := c.extractOptionalFloat(result, "askKraken")
		if !ok {
			continue
		}

		bidBitsamp, ok := c.extractOptionalFloat(result, "bidBitstamp")
		if !ok {
			continue
		}

		askBitstamp, ok := c.extractOptionalFloat(result, "askBitstamp")
		if !ok {
			continue
		}

		symbolData := symbol.Symbol{
			Pair:          pair,
			Price:         priceLast,
			Price1hAgo:    price1hAgo,
			Price1dAgo:    price1dAgo,
			Price7dAgo:    price7dAgo,
			PriceChange1h: priceChange1h,
			PriceChange1d: priceChange1d,
			PriceChange7d: priceChange7d,
			PriceBitstamp: priceBitstamp,
			PriceKraken:   priceKraken,
			BidKraken:     bidKraken,
			AskKraken:     askKraken,
			BidBitstamp:   bidBitsamp,
			AskBitstamp:   askBitstamp,
		}

		c.calculateArbitrage(ctx, pipe, &symbolData)

		// Convert symbolData to json
		symbolJsonData, err := json.Marshal(symbolData)
		if err != nil {
			c.appContext.Logger.Printf("Failed to marshal priceData to JSON: %v \n", err)
			continue
		}
		// Add data to the pipeline. It will not be executed until pipe.Exec() is called.
		pipe.Set(ctx, "hash:"+pair, symbolJsonData, 0)
		if symbolData.Arbitrage != nil {
			pipe.ZAdd(ctx, "zset:arbitrage", redis.Z{Member: symbolData.Pair, Score: *symbolData.Arbitrage})
		}

		task := workerpool.Task{
			Execute: func() error {
				// fmt.Printf("Der preis von %s hat sich geändert von %s nach %2f \n", pair, previousPrice, priceLast)
				return c.updateSymbolPriceChanges(&symbolJsonData)
			},
			Pair: pair,
		}

		// Übermittle die Task an den WorkerPool
		c.pool.Tasks <- task

	}

	if result.Err() != nil {
		c.appContext.Logger.Printf("Failed to iterate through the query result from fetchSymbolPriceChanges: %v \n", err)
	}

	// Execute all buffered commands at once.
	_, err = pipe.Exec(ctx)
	if err != nil {
		c.appContext.Logger.Printf("Could not execute pipeline from fetchSymbolPriceChanges: %v \n", err)
	}
}

// extractFloatValue attempts to extract the value for a given key from a Record
// as a float64. On success, it returns the extracted value and true.
// In case of failure, an error message is logged, and it returns 0 and false.func (c *JobController) extractFloatValue(result *api.QueryTableResult, key string) (float64, bool) {
func (c *JobController) extractFloatValue(result *api.QueryTableResult, key string) (float64, bool) {
	rawValue := result.Record().ValueByKey(key)
	value, ok := rawValue.(float64)
	if !ok {
		c.appContext.Logger.Printf("Error: Value %s is not a float64\n", key)
	}
	return value, ok
}

func (c *JobController) extractOptionalFloat(result *api.QueryTableResult, key string) (*float64, bool) {
	rawValue := result.Record().ValueByKey(key)
	if rawValue == nil {
		return nil, true
	}

	floatVal, ok := rawValue.(float64)
	if !ok {
		c.appContext.Logger.Printf("Error: Value %s could not be converted to float64\n", key)
		return nil, false
	}
	return &floatVal, true
}

// updateSymbolPriceChanges updates the symbol price changes in Redis.
func (c *JobController) updateSymbolPriceChanges(jsonData *[]byte) error {
	// Senden als POST Request
	resp, err := http.Post(c.wsURL, "application/json", bytes.NewBuffer(*jsonData))
	if err != nil {
		return fmt.Errorf("error to send post request err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status code %d", resp.StatusCode)
	}

	return nil
}
