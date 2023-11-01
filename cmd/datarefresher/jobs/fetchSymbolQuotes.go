package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c *JobController) FetchSymbolQuotes() {
	// Logic to fetch symbols and data from InfluxDB

	// Use a.appContext.QueryInfluxDB to fetch the data from InfluxDB
	result, err := c.appContext.InfluxClient.Query("queries/getSymbolQuotes.flux")
	if err != nil {
		c.appContext.Logger.Fatalf("Failed to query get_symbols: %s", err)
	}

	// Create a context.Context with a timeout of 100 milliseconds
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	// Make sure to clean up the context at the end
	defer cancel()

	// Create a pipeline for redis HSet USD
	pipe := c.appContext.Redis.Pipeline()

	// Iterate over the QueryTableResult and extract the data
	for result.Next() {
		pair, ok := result.Record().ValueByKey("pair").(string)
		if !ok {
			c.appContext.Logger.Println("Error: pair is not a string")
			continue
		}
		qv, ok := result.Record().ValueByKey("quote_volume").(float64)
		if !ok {
			c.appContext.Logger.Println("Error: quote_volume is not a float64")
			continue
		}

		var key string
		// Check if the pair has "/USD" or "/EUR" and set the Redis key accordingly.
		if strings.HasSuffix(pair, "/USD") {
			key = "zset:symbols:usd"
			// pair = strings.TrimSuffix(pair, "/USD")
		} else if strings.HasSuffix(pair, "/EUR") {
			key = "zset:symbols:eur"
			// pair = strings.TrimSuffix(pair, "/EUR")
		} else {
			c.appContext.Logger.Println("Error: unexpected pair suffix")
			continue
		}
		pipe.ZAdd(ctx, key, redis.Z{Member: pair, Score: qv})

	}

	if result.Err() != nil {
		c.appContext.Logger.Printf("Failed to iterate through the query result from get_symbols: %v", err)
	}

	// Execute all buffered commands at once.
	_, err = pipe.Exec(ctx)
	if err != nil {
		c.appContext.Logger.Fatalf("Could not execute pipeline for fetch Symbol Quotes: %v", err)
	}
}
