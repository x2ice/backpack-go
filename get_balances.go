package backpack

import (
	"context"
	"fmt"
	"net/http"
)

type Balances map[string]struct {
	Available float64 `json:"available,string"`
	Locked    float64 `json:"locked,string"`
	Staked    float64 `json:"staked,string"`
}

func (c *Backpack) GetBalances(ctx context.Context) (Balances, error) {
	url := fmt.Sprintf("%s/api/v1/capital", HTTP_API)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header, err = c.Headers("balanceQuery", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	balances, err := decodeResponse[Balances](resp.Body)
	return *balances, err
}
