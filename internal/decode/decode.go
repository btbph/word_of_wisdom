package decode

import (
	"encoding/json"
	"io"
)

func JsonFromReader[T any](data io.Reader) (T, error) {
	var req T
	if err := json.NewDecoder(data).Decode(&req); err != nil {
		return req, nil
	}

	return req, nil
}
