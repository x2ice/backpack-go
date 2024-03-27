package backpack

import (
	"bytes"
	"context"
	"fmt"
)

type DepthUpdate struct {
	EventType       string     `json:"e"`
	EventTimestamp  int64      `json:"E"`
	EngineTimestamp int64      `json:"T"`
	FirstUpdateID   int64      `json:"U"`
	FinalUpdateID   int64      `json:"u"`
	Asks            [][]string `json:"a"`
	Bids            [][]string `json:"b"`
	Symbol          string     `json:"s"`
	Stream          string     `json:"stream"`
}

func (c *Backpack) SubscribeOrderDepth(ctx context.Context, symbol string) *Subscription[DepthUpdate] {
	sub := new(Subscription[DepthUpdate])

	sub.Msgs = make(chan *DepthUpdate)
	sub.unsubscribe = make(chan struct{})

	msgs := make(chan []byte)

	stream := fmt.Sprintf("depth.%s", symbol)
	c.Subscribe(ctx, stream, false, msgs, sub.unsubscribe)

	go func() {
		defer close(sub.Msgs)
		defer close(sub.unsubscribe)

		for msg := range msgs {
			reader := bytes.NewReader(msg[8:])

			depthUpdate, err := decodeResponse[DepthUpdate](reader)
			if err != nil {
				sub.Err <- err
				continue
			}

			sub.Msgs <- depthUpdate
		}
	}()

	return sub
}
