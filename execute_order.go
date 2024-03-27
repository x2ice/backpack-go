package backpack

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	json "github.com/goccy/go-json"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Order struct {
	OrderType             string `json:"orderType"`
	ID                    string `json:"id"`
	ClientID              int    `json:"clientId,omitempty"`
	Symbol                string `json:"symbol"`
	Side                  string `json:"side"`
	Quantity              string `json:"quantity,omitempty"`
	ExecutedQuantity      string `json:"executedQuantity"`
	QuoteQuantity         string `json:"quoteQuantity"`
	ExecutedQuoteQuantity string `json:"executedQuoteQuantity"`
	TriggerPrice          string `json:"triggerPrice,omitempty"`
	TimeInForce           string `json:"timeInForce"`
	SelfTradePrevention   string `json:"selfTradePrevention"`
	Status                string `json:"status"`
	CreatedAt             int64  `json:"createdAt"`
}

func (c *Backpack) ExecuteOrder(ctx context.Context, orderType, side, symbol string, options map[string]any) (*Order, error) {
	options["orderType"] = cases.Title(language.English).String(orderType)

	side = cases.Title(language.English).String(side)
	if side != "Ask" && side != "Bid" {
		return nil, ErrInvalidOrderSide
	}

	options["side"] = side
	options["symbol"] = strings.ToUpper(symbol)

	payload, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/order", HTTP_API)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header, err = c.Headers("orderExecute", options)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		switch string(b) {
		case "Insufficient funds":
			return nil, err
		}
	}

	return decodeResponse[Order](resp.Body)
}
