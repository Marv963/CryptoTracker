package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/Marv963/CryptoTracker/app/pkg/appcontext"
)

type SymbolController struct {
	appContext *appcontext.AppContext
}

func NewSymbolController(appContext *appcontext.AppContext) *SymbolController {
	return &SymbolController{
		appContext: appContext,
	}
}

func (s *SymbolController) GetSymbols(c *fiber.Ctx) error {
	// Logik zum Abrufen der Symbole und Daten von InfluxDB
	page, err := c.ParamsInt("page", 1)
	if err != nil {
		return c.Status(400).SendString("Bad Request")
	}
	currency := c.Query("currency", "usd")

	// Berechnen des Start- und End-Indexes für die Abfrage
	startIndex := (page - 1) * 100
	endIndex := startIndex + 99

	results, err := s.getSortedSet(int(startIndex), int(endIndex), "zset:symbols:"+currency)
	if err != nil {
		s.appContext.Logger.Printf("%v", err)
		return c.Status(500).SendString("Internal Server Error")
	}
	hashes, err := s.getHashes(results)
	if err != nil {
		s.appContext.Logger.Printf("%v", err)
		return c.Status(500).SendString("Internal Server Error")
	}
	// Setzen des Content-Type Headers auf application/json
	// c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.JSON(hashes)

}

func (s *SymbolController) getSortedSet(startIndex int, endIndex int, hash string) ([]redis.Z, error) {

	// Create a context.Context with a timeout of 100 milliseconds
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	// Make sure to clean up the context at the end
	defer cancel()

	// Abfrage der Top 100 Symbole aus Redis.
	results, err := s.appContext.Redis.ZRevRangeWithScores(ctx, hash, int64(startIndex), int64(endIndex)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get top 100 from sorted set %s err: %v", hash, err)
	}
	return results, nil
}

func (s *SymbolController) getHashes(results []redis.Z) ([]map[string]interface{}, error) {
	// Create a context.Context with a timeout of 100 milliseconds
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	// Make sure to clean up the context at the end
	defer cancel()

	// Pipeline erstellen um eine weitere Abfrage an Redis zu senden
	pipe := s.appContext.Redis.Pipeline()

	// Create a slice to hold your data points
	for _, result := range results {

		pipe.Get(ctx, "hash:"+result.Member.(string))
	}

	// Ausführen aller Befehle in der Pipeline
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash from symbols err: %v", err)
	}
	data := make([]map[string]interface{}, 0, len(cmds))
	for _, cmd := range cmds {
		val, err := cmd.(*redis.StringCmd).Bytes()
		if err != nil {
			s.appContext.Logger.Printf("Failed to get value from redis command: %v \n", err)
			continue
		}

		var symbolData map[string]interface{}
		if err := json.Unmarshal(val, &symbolData); err != nil {
			s.appContext.Logger.Printf("Failed to unmarshal json data: %v \n", err)
			continue
		}

		data = append(data, symbolData)
	}
	return data, nil
}
