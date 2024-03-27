package backpack

import (
	"crypto/ed25519"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"
)

var (
	HTTP_API = "https://api.backpack.exchange"
	WS_API   = "wss://ws.backpack.exchange"

	WebsocketTimeout   = time.Second * 60
	WebsocketKeepAlive = true
)

type Backpack struct {
	HttpClient      *http.Client
	privateKey      ed25519.PrivateKey
	VerifyingKeyB64 string
}

func NewBackpack() *Backpack {
	return &Backpack{
		HttpClient: &http.Client{},
	}
}

func (c *Backpack) SetAPISecret(apiSecret string) *Backpack {
	apiSecretB64, err := base64.StdEncoding.DecodeString(apiSecret)
	if err != nil {
		panic(err)
	}

	privateKey := ed25519.NewKeyFromSeed(apiSecretB64)
	c.privateKey = privateKey

	verifyingKey := privateKey.Public().(ed25519.PublicKey)
	verifyingKeyB64 := base64.StdEncoding.EncodeToString(verifyingKey)
	c.VerifyingKeyB64 = verifyingKeyB64

	return c
}

func (c *Backpack) SetProxy(proxy string) *Backpack {
	url, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}

	transport := &http.Transport{Proxy: http.ProxyURL(url)}
	c.HttpClient.Transport = transport

	return c
}
