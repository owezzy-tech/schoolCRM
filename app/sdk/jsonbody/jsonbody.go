// Package jsonbody strictly decodes size-limited JSON request bodies.
package jsonbody

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Decode reads at most limit bytes and decodes exactly one JSON value into v,
// rejecting unknown fields.
func Decode(r *http.Request, v any, limit int64) error {
	data, err := io.ReadAll(io.LimitReader(r.Body, limit+1))
	if err != nil {
		return err
	}
	if int64(len(data)) > limit {
		return fmt.Errorf("request exceeds %d KiB", limit/1024)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return errors.New("expected a single JSON object")
	}
	return nil
}
