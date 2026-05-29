package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"go.opentelemetry.io/otel/attribute"
)

const jsonAPIContentType = "application/vnd.api+json"

type jsonAPIDocument struct {
	JSONAPI jsonAPIObject  `json:"jsonapi"`
	Data    any            `json:"data,omitempty"`
	Errors  []jsonAPIError `json:"errors,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

type jsonAPIObject struct {
	Version string `json:"version"`
}

type jsonAPIResource struct {
	Type       string `json:"type"`
	ID         string `json:"id,omitempty"`
	Attributes any    `json:"attributes,omitempty"`
}

type jsonAPIError struct {
	Status string         `json:"status"`
	Code   string         `json:"code,omitempty"`
	Title  string         `json:"title"`
	Detail string         `json:"detail,omitempty"`
	Meta   map[string]any `json:"meta,omitempty"`
}

type queryResultEnvelope struct {
	Items       json.RawMessage `json:"items"`
	Total       int             `json:"total"`
	Page        int             `json:"page"`
	RowsPerPage int             `json:"rowsPerPage"`
}

type encodedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NoResponse tells the Respond function to not respond to the request. In these
// cases the app layer code has already done so.
type NoResponse struct{}

// NewNoResponse constructs a no response value.
func NewNoResponse() NoResponse {
	return NoResponse{}
}

// Encode implements the Encoder interface.
func (NoResponse) Encode() ([]byte, string, error) {
	return nil, "", nil
}

// =============================================================================

type httpStatus interface {
	HTTPStatus() int
}

// Respond sends a response to the client.
func Respond(ctx context.Context, w http.ResponseWriter, resp Encoder) error {
	if _, ok := resp.(NoResponse); ok {
		return nil
	}

	// If the context has been canceled, it means the client is no longer
	// waiting for a response.
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("client disconnected, do not send response")
		}
	}

	statusCode := http.StatusOK

	switch v := resp.(type) {
	case httpStatus:
		statusCode = v.HTTPStatus()

	case error:
		statusCode = http.StatusInternalServerError

	default:
		if resp == nil {
			statusCode = http.StatusNoContent
		}
	}

	_, span := addSpan(ctx, "web.send.response", attribute.Int("status", statusCode))
	defer span.End()

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	data, contentType, err := encodeJSONAPI(resp, statusCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("respond: encode: %w", err)
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("respond: write: %w", err)
	}

	return nil
}

func encodeJSONAPI(resp Encoder, statusCode int) ([]byte, string, error) {
	data, _, err := resp.Encode()
	if err != nil {
		return nil, "", err
	}

	if _, ok := resp.(error); ok {
		var appErr encodedError
		if err := json.Unmarshal(data, &appErr); err != nil {
			appErr.Message = resp.(error).Error()
		}

		doc := jsonAPIDocument{
			JSONAPI: jsonAPIObject{Version: "1.1"},
			Errors: []jsonAPIError{
				{
					Status: fmt.Sprintf("%d", statusCode),
					Code:   appErr.Code,
					Title:  http.StatusText(statusCode),
					Detail: appErr.Message,
				},
			},
		}

		data, err := json.Marshal(doc)
		return data, jsonAPIContentType, err
	}

	doc := jsonAPIDocument{JSONAPI: jsonAPIObject{Version: "1.1"}}
	if query, ok := decodeQueryResult(data); ok {
		resources, err := collectionResources(query.Items, resourceType(resp))
		if err != nil {
			return nil, "", err
		}

		doc.Data = resources
		doc.Meta = map[string]any{
			"total":       query.Total,
			"page":        query.Page,
			"rowsPerPage": query.RowsPerPage,
		}
	} else {
		doc.Data = jsonAPIResource{
			Type:       resourceType(resp),
			ID:         resourceID(data),
			Attributes: resourceAttributes(data),
		}
	}

	encoded, err := json.Marshal(doc)
	return encoded, jsonAPIContentType, err
}

func decodeQueryResult(data []byte) (queryResultEnvelope, bool) {
	var query queryResultEnvelope
	if err := json.Unmarshal(data, &query); err != nil {
		return queryResultEnvelope{}, false
	}

	if query.Items == nil || query.Page == 0 || query.RowsPerPage == 0 {
		return queryResultEnvelope{}, false
	}

	return query, true
}

func collectionResources(items json.RawMessage, typ string) ([]jsonAPIResource, error) {
	var rawItems []json.RawMessage
	if err := json.Unmarshal(items, &rawItems); err != nil {
		return nil, fmt.Errorf("decode collection items: %w", err)
	}

	resources := make([]jsonAPIResource, 0, len(rawItems))
	for _, item := range rawItems {
		resources = append(resources, jsonAPIResource{
			Type:       typ,
			ID:         resourceID(item),
			Attributes: resourceAttributes(item),
		})
	}

	return resources, nil
}

func resourceType(resp Encoder) string {
	typ := reflect.TypeOf(resp)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	typeName := typ.Name()
	if typeName == "" {
		typeName = "resource"
	}

	return normalizeResourceType(typeName)
}

func normalizeResourceType(typeName string) string {
	if typeName == "" {
		return "resource"
	}

	if open := strings.LastIndex(typeName, "["); open >= 0 {
		inner := strings.TrimSuffix(typeName[open+1:], "]")
		if dot := strings.LastIndex(inner, "."); dot >= 0 {
			inner = inner[dot+1:]
		}

		return strings.ToLower(inner)
	}

	return strings.ToLower(typeName)
}

func resourceID(data []byte) string {
	var body map[string]json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return ""
	}

	var id string
	if rawID, ok := body["id"]; ok {
		if err := json.Unmarshal(rawID, &id); err == nil {
			return id
		}
	}

	return ""
}

func resourceAttributes(data []byte) json.RawMessage {
	var body map[string]json.RawMessage
	if err := json.Unmarshal(data, &body); err != nil {
		return json.RawMessage(data)
	}

	delete(body, "id")

	encoded, err := json.Marshal(body)
	if err != nil {
		return json.RawMessage(data)
	}

	return encoded
}
