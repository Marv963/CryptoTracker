package ccxt

type Info struct {
	Open            string `json:"open"`
	High            string `json:"high"`
	Low             string `json:"low"`
	Last            string `json:"last"`
	Volume          string `json:"volume"`
	Vwap            string `json:"vwap"`
	Bid             string `json:"bid"`
	Ask             string `json:"ask"`
	Side            string `json:"side"`
	Open24          string `json:"open_24"`
	PercentChange24 string `json:"percent_change_24"`
}

type Message struct {
	Info          Info     `json:"info"`
	BidVolume     *float64 `json:"bidVolume"`
	AskVolume     *float64 `json:"askVolume"`
	PreviousClose *float64 `json:"previousClose"`
	Symbol        string   `json:"symbol"`
	Datetime      string   `json:"datetime"`
	Timestamp     int64    `json:"timestamp"`
	High          float64  `json:"high"`
	Low           float64  `json:"low"`
	Bid           float64  `json:"bid"`
	Ask           float64  `json:"ask"`
	Vwap          float64  `json:"vwap"`
	Open          float64  `json:"open"`
	Close         float64  `json:"close"`
	Last          float64  `json:"last"`
	Change        float64  `json:"change"`
	Percentage    float64  `json:"percentage"`
	Average       float64  `json:"average"`
	BaseVolume    float64  `json:"baseVolume"`
	QuoteVolume   float64  `json:"quoteVolume"`
}
