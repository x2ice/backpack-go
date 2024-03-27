package backpack

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type AuthenticationMessage struct {
	VerifyingKeyB64 string
	Timestamp       string
	Window          string
	SignatureB64    string
}

func (c *Backpack) SignRequest(instruction string, args map[string]any) (*AuthenticationMessage, error) {
	if c.privateKey == nil {
		return nil, ErrCredentialsRequired
	}

	message := fmt.Sprintf("instruction=%s", instruction)
	if args != nil {
		keys := make([]string, 0, len(args))
		for k := range args {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for _, k := range keys {
			v := args[k].(string)
			message += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(v)
		}
	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	message += fmt.Sprintf("&timestamp=%s", timestamp)

	window := "5000"
	message += fmt.Sprintf("&window=%s", window)

	signature := ed25519.Sign(c.privateKey, []byte(message))
	signatureB64 := base64.StdEncoding.EncodeToString(signature)

	return &AuthenticationMessage{
		c.VerifyingKeyB64,
		timestamp,
		window,
		signatureB64,
	}, nil
}

func (c *Backpack) Headers(instruction string, args map[string]any) (http.Header, error) {
	auth, err := c.SignRequest(instruction, args)
	if err != nil {
		return nil, err
	}

	return http.Header{
		"X-API-KEY":    []string{auth.VerifyingKeyB64},
		"X-TIMESTAMP":  []string{auth.Timestamp},
		"X-WINDOW":     []string{auth.Window},
		"Content-Type": []string{"application/json"},
		"X-SIGNATURE":  []string{auth.SignatureB64},
	}, nil

}
