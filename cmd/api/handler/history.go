package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Marv963/CryptoTracker/app/pkg/model/symbol"
)

func (s *SymbolController) GetSymbolHistory(c *fiber.Ctx) error {
	pair := c.Params("pair")
	pair = strings.Replace(pair, "-", "/", 1)

	result, err := s.appContext.InfluxClient.Query("queries/getSymbolHistory.flux", pair)
	if err != nil {
		s.appContext.Logger.Printf("Failed to query symbol history for pair %s: %v\n", pair, err)
		return c.Status(500).SendString("Internal Server Error")
	}

	// Slice to hold our history entries.
	dataMap := make(map[string][]symbol.SymbolHistory)

	for result.Next() {
		value, ok := result.Record().ValueByKey("_value").(float64)
		if !ok {
			s.appContext.Logger.Println("Error: value is not a float64")
			return c.Status(500).SendString("Internal Server Error")
		}
		time, ok := result.Record().ValueByKey("_time").(int64)
		if !ok {
			s.appContext.Logger.Println("Error: value is not a int")
			return c.Status(500).SendString("Internal Server Error")
		}
		exchange, ok := result.Record().ValueByKey("exchange").(string)
		if !ok {
			s.appContext.Logger.Println("Error: exchange is not a string")
			return c.Status(500).SendString("Internal Server Error")
		}

		// Add data to map, organized by exchange.
		dataMap[exchange] = append(dataMap[exchange], symbol.SymbolHistory{
			Value: value,
			Time:  time / 1000000000, // InfluxDB unix time is Nano seconds, but seconds are needed
		})
	}
	formattedData := []symbol.SymbolHistoryExchangeData{}

	// Convert the map data into the desired format.
	for exchange, data := range dataMap {
		formattedData = append(formattedData, symbol.SymbolHistoryExchangeData{
			Exchange: exchange,
			Data:     data,
		})
	}

	// Return the history slice as JSON.
	return c.JSON(formattedData)
}
