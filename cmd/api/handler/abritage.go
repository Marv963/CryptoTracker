package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (s *SymbolController) GetArbitrage(c *fiber.Ctx) error {
	results, err := s.getSortedSet(0, 100, "zset:arbitrage")
	if err != nil {
		s.appContext.Logger.Printf("%v", err)
		return c.Status(500).SendString("Internal Server Error")
	}
	hashes, err := s.getHashes(results)
	if err != nil {
		s.appContext.Logger.Printf("%v", err)
		return c.Status(500).SendString("Internal Server Error")
	}
	return c.JSON(hashes)
}
