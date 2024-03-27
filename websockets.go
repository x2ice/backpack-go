package backpack

import (
	"context"
	"log"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

func keepAlive(conn *websocket.Conn, timeout time.Duration, quit chan struct{}) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-quit:
				return
			default:
				deadline := time.Now().Add(10 * time.Second)
				err := conn.WriteControl(websocket.PingMessage, []byte{}, deadline)
				if err != nil {
					return
				}

				<-ticker.C
				if time.Since(lastResponse) > timeout {
					conn.Close()
					return
				}
			}
		}
	}()
}

type Subscription[T any] struct {
	Msgs        chan *T
	Err         chan error
	unsubscribe chan struct{}
}

func (c *Backpack) unsubscribe(conn *websocket.Conn, stream string) {
	payload := map[string]any{
		"method": "UNSUBSCRIBE",
		"params": []string{stream},
	}

	msg, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	conn.WriteMessage(websocket.TextMessage, msg)
}

func (s *Subscription[T]) Unsubscribe() {
	s.unsubscribe <- struct{}{}
}

func (c *Backpack) Subscribe(ctx context.Context, stream string, signRequired bool, msgs chan []byte, unsubscribe chan struct{}) {
	conn, _, err := websocket.DefaultDialer.Dial(WS_API, nil)
	if err != nil {
		panic(err)
	}

	go func() {

		defer func() {
			c.unsubscribe(conn, stream)
			close(msgs)
			conn.Close()
		}()

		payload := map[string]any{
			"method": "SUBSCRIBE",
			"params": []string{stream},
		}

		if signRequired || c.privateKey != nil {
			instruction := "subscribe"
			auth, err := c.SignRequest(instruction, nil)
			if err != nil {
				panic(err)
			}

			payload["signature"] = []string{
				auth.VerifyingKeyB64,
				auth.SignatureB64,
				auth.Timestamp,
				auth.Window,
			}
		}

		msg, _ := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			panic(err)
		}

		if WebsocketKeepAlive {
			keepAlive(conn, WebsocketTimeout, unsubscribe)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-unsubscribe:
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					panic(err)
				}
				msgs <- msg
			}
		}
	}()
}
