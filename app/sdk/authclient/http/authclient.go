// Package http provides support to access the auth service.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"time"

	"github.com/owezzy/schoolCRM/app/sdk/authclient"
	"github.com/owezzy/schoolCRM/app/sdk/errs"
	"github.com/owezzy/schoolCRM/foundation/logger"
	"github.com/owezzy/schoolCRM/foundation/otel"
	"go.opentelemetry.io/otel/attribute"
)

// This provides a default client configuration, but it's recommended
// this is replaced by the user with application specific settings using
// the WithClient function at the time a AuthAPI is constructed.
// DualStack Deprecated: Fast Fallback is enabled by default. To disable, set FallbackDelay to a negative value.
var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

// Client represents a client that can talk to the auth service.
type Client struct {
	log  *logger.Logger
	url  string
	http *http.Client
}

type jsonAPIDocument struct {
	Data   jsonAPIResource `json:"data"`
	Errors []jsonAPIError `json:"errors"`
}

type jsonAPIResource struct {
	Attributes json.RawMessage `json:"attributes"`
}

type jsonAPIError struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// New constructs an Auth that can be used to talk with the auth service.
func New(log *logger.Logger, url string, options ...func(cln *Client)) (*Client, error) {
	cln := Client{
		log:  log,
		url:  url,
		http: &defaultClient,
	}

	for _, option := range options {
		option(&cln)
	}

	return &cln, nil
}

// WithClient adds a custom client for processing requests. It's recommend
// to not use the default client and provide your own.
func WithClient(http *http.Client) func(cln *Client) {
	return func(cln *Client) {
		cln.http = http
	}
}

// Authenticate calls the auth service to authenticate the user.
func (cln *Client) Authenticate(ctx context.Context, authorization string) (authclient.AuthenticateResp, error) {
	endpoint := fmt.Sprintf("%s/v1/auth/authenticate", cln.url)

	headers := map[string]string{
		"authorization": authorization,
	}

	var resp authclient.AuthenticateResp
	if err := cln.do(ctx, http.MethodGet, endpoint, headers, nil, &resp); err != nil {
		return authclient.AuthenticateResp{}, err
	}

	return resp, nil
}

// Authorize calls the auth service to authorize the user.
func (cln *Client) Authorize(ctx context.Context, auth authclient.Authorize) error {
	endpoint := fmt.Sprintf("%s/v1/auth/authorize", cln.url)

	if err := cln.do(ctx, http.MethodPost, endpoint, nil, auth, nil); err != nil {
		return err
	}

	return nil
}

func (cln *Client) do(ctx context.Context, method string, endpoint string, headers map[string]string, body any, v any) error {
	var statusCode int

	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parsing endpoint: %w", err)
	}
	base := path.Base(u.Path)

	cln.log.Info(ctx, "authclient: rawRequest: started", "method", method, "call", base, "endpoint", endpoint)
	defer func() {
		cln.log.Info(ctx, "authclient: rawRequest: completed", "status", statusCode)
	}()

	ctx, span := otel.AddSpan(ctx, fmt.Sprintf("app.sdk.authclient.%s", base), attribute.String("endpoint", endpoint))
	defer func() {
		span.SetAttributes(attribute.Int("status", statusCode))
		span.End()
	}()

	var b bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&b).Encode(body); err != nil {
			return fmt.Errorf("encoding error: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, &b)
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for key, value := range headers {
		cln.log.Info(ctx, "authclient: rawRequest", "key", key, "value", value)
		req.Header.Set(key, value)
	}

	otel.AddTraceToRequest(ctx, req)

	resp, err := cln.http.Do(req)
	if err != nil {
		return fmt.Errorf("do: error: %w", err)
	}
	defer resp.Body.Close()

	// Assign so it can be logged in the defer above.
	statusCode = resp.StatusCode

	if statusCode == http.StatusNoContent {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	switch statusCode {
	case http.StatusOK:
		if err := decodeSuccess(data, v); err != nil {
			return fmt.Errorf("failed: response: %s, decoding error: %w ", string(data), err)
		}
		return nil

	case http.StatusUnauthorized, http.StatusForbidden:
		err, decodeErr := decodeError(data)
		if decodeErr != nil {
			return fmt.Errorf("failed: response: %s, decoding error: %w ", string(data), decodeErr)
		}
		return err

	default:
		return fmt.Errorf("failed: response: %s", string(data))
	}
}

func decodeSuccess(data []byte, v any) error {
	if isNilDecodeTarget(v) {
		return nil
	}

	var doc jsonAPIDocument
	if err := json.Unmarshal(data, &doc); err == nil && len(doc.Data.Attributes) > 0 {
		return json.Unmarshal(doc.Data.Attributes, v)
	}

	return json.Unmarshal(data, v)
}

func decodeError(data []byte) (*errs.Error, error) {
	var doc jsonAPIDocument
	if err := json.Unmarshal(data, &doc); err == nil && len(doc.Errors) > 0 {
		apiErr := doc.Errors[0]
		code := errorCode(apiErr.Code)
		message := apiErr.Detail
		if message == "" {
			message = apiErr.Title
		}
		if message == "" {
			message = code.String()
		}

		return errs.Errorf(code, "%s", message), nil
	}

	var appErr errs.Error
	if err := json.Unmarshal(data, &appErr); err != nil {
		return nil, err
	}

	return &appErr, nil
}

func errorCode(code string) errs.ErrCode {
	switch code {
	case errs.PermissionDenied.String():
		return errs.PermissionDenied
	case errs.Unauthenticated.String():
		return errs.Unauthenticated
	default:
		return errs.Unauthenticated
	}
}

func isNilDecodeTarget(v any) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	return value.Kind() == reflect.Pointer && value.IsNil()
}

func (cln *Client) Close() error {
	return nil
}
