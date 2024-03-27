package backpack

import (
	"io"

	json "github.com/goccy/go-json"
)

func decodeResponse[V any](r io.Reader) (*V, error) {
	var v V
	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
