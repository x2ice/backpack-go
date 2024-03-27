package backpack

import (
	"bytes"
	"context"
)

type OrderUpdate struct {
	EventType           string  `json:"e"`
	EventTime           int64   `json:"E"`
	Symbol              string  `json:"s"`
	ClientOrderID       string  `json:"c"`
	Side                string  `json:"S"`
	OrderType           string  `json:"o"`
	TimeInForce         string  `json:"f"`
	Quantity            float64 `json:"q,string"`
	QuantityInQuote     float64 `json:"Q,string"`
	Price               float64 `json:"p,string"`
	TriggerPrice        string  `json:"P"`
	OrderState          string  `json:"X"`
	OrderID             string  `json:"i"`
	FillQuantity        string  `json:"l"`
	ExecutedQuantity    string  `json:"z"`
	ExecutedInQuote     string  `json:"Z"`
	FillPrice           string  `json:"L"`
	IsMaker             bool    `json:"m"`
	Fee                 float64 `json:"n,string"`
	FeeSymbol           string  `json:"N"`
	SelfTradePrevention string  `json:"V"`
	EngineTimestamp     int64   `json:"T"`
}

func (c *Backpack) SubscribeOrderUpdate(ctx context.Context) *Subscription[OrderUpdate] {
	sub := new(Subscription[OrderUpdate])

	sub.Msgs = make(chan *OrderUpdate)
	sub.unsubscribe = make(chan struct{})
	msgs := make(chan []byte)

	c.Subscribe(ctx, "account.orderUpdate", true, msgs, sub.unsubscribe)

	go func() {
		defer close(sub.Msgs)
		defer close(sub.unsubscribe)

		for msg := range msgs {
			reader := bytes.NewReader(msg[8:])

			orderUpdate, err := decodeResponse[OrderUpdate](reader)
			if err != nil {
				sub.Err <- err
				continue
			}

			sub.Msgs <- orderUpdate
		}
	}()

	return sub
}
