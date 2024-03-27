package backpack

import (
	"context"
	"fmt"
	"net/http"
)

type Market struct {
	Symbol      string `json:"symbol"`
	BaseSymbol  string `json:"baseSymbol"`
	QuoteSymbol string `json:"quoteSymbol"`
	Filters     *struct {
		Leverage string `json:"leverage,omitempty"`
		Price    *struct {
			Max      float64 `json:"maxPrice,,string,omitempty"`
			Min      float64 `json:"minPrice,string"`
			TickSize float64 `json:"tickSize,string"`
		} `json:"price"`
		Quantity *struct {
			Max      float64 `json:"maxQuantity,string,omitempty"`
			Min      float64 `json:"minQuantity,string"`
			StepSize float64 `json:"stepSize,string"`
		} `json:"quantity"`
	} `json:"filters"`
}

type Markets []Market

func (c *Backpack) GetMarkets(ctx context.Context) (Markets, error) {
	url := fmt.Sprintf("%s/api/v1/markets", HTTP_API)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	markets, err := decodeResponse[Markets](resp.Body)
	return *markets, err
}
