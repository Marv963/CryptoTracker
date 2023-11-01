package symbol

type Symbol struct {
	PriceBitstamp        *float64 `json:"priceBitstamp"`
	PriceKraken          *float64 `json:"priceKraken"`
	LowestPrice          *float64 `json:"lowestPrice"`
	HighestPrice         *float64 `json:"highestPrice"`
	Arbitrage            *float64 `json:"arbitrage"`
	BidKraken            *float64 `json:"bidKraken"`
	AskKraken            *float64 `json:"askKraken"`
	BidBitstamp          *float64 `json:"bidBitstamp"`
	AskBitstamp          *float64 `json:"askBitstamp"`
	LowestPriceExchange  string   `json:"lowestPriceExchange"`
	HighestPriceExchange string   `json:"highestPriceExchange"`
	Pair                 string   `json:"pair"`
	Price                float64  `json:"price"`
	Price1hAgo           float64  `json:"price1hAgo"`
	Price1dAgo           float64  `json:"price1dAgo"`
	Price7dAgo           float64  `json:"price7dAgo"`
	PriceChange1h        float64  `json:"priceChange1h"`
	PriceChange1d        float64  `json:"priceChange1d"`
	PriceChange7d        float64  `json:"priceChange7d"`
	QuoteVolume          float64  `json:"quoteVolume"`
}

type SymbolQuote struct {
	CC          string  `json:"cc"`
	QuoteVolume float64 `json:"quote_volume"`
}

type SymbolHistory struct {
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
}

type SymbolHistoryExchangeData struct {
	Exchange string          `json:"exchange"`
	Data     []SymbolHistory `json:"data"`
}
