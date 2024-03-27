package backpack

import (
	"context"
	"fmt"
	"net/http"
)

type DepositAddress struct {
	Address string
}

func (c *Backpack) GetDepositAddress(ctx context.Context, blockchain string) (*DepositAddress, error) {
	url := fmt.Sprintf("%s/wapi/v1/capital/deposit/address?blockchain=%s", HTTP_API, blockchain)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	params := map[string]any{"blockchain": string(blockchain)}
	req.Header, err = c.Headers("depositAddressQuery", params)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeResponse[DepositAddress](resp.Body)
}
