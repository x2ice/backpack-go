package backpack

import (
	"bytes"
	"context"
	"fmt"
)

type KLine struct {
	EventType      string  `json:"e"`
	EventTimestamp int64   `json:"E"`
	ClosePrice     string  `json:"c,string"`
	OpenPrice      string  `json:"o,string"`
	HighPrice      float64 `json:"h,string"`
	LowPrice       float64 `json:"l,string"`
	NumberOfTrades int64   `json:"n"`
	Symbol         string  `json:"s"`
	Time           string  `json:"t"`
	Volume         float64 `json:"v,string"`
	IsClosed       bool    `json:"X"`
	Stream         string  `json:"stream"`
}

func (c *Backpack) SubscribeKLine(ctx context.Context, interval, symbol string) *Subscription[KLine] {
	sub := new(Subscription[KLine])

	sub.Msgs = make(chan *KLine)
	msgs := make(chan []byte)

	stream := fmt.Sprintf("kline.%s.%s", interval, symbol)
	c.Subscribe(ctx, stream, false, msgs, sub.unsubscribe)

	go func() {
		defer close(sub.Msgs)
		defer close(sub.unsubscribe)

		for msg := range msgs {
			reader := bytes.NewReader(msg[8:])

			kline, err := decodeResponse[KLine](reader)
			if err != nil {
				sub.Err <- err
				continue
			}

			sub.Msgs <- kline
		}
	}()

	return sub
}
