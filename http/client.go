package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/gsmcwhirter/go-util/v12/deferutil"
	"github.com/gsmcwhirter/go-util/v12/errors"
	"github.com/gsmcwhirter/go-util/v12/json"
	"github.com/gsmcwhirter/go-util/v12/logging"
	"github.com/gsmcwhirter/go-util/v12/telemetry"
)

type (
	Request             = retryablehttp.Request
	ResponseHandlerFunc = retryablehttp.ResponseHandlerFunc
)

type RetryOptions struct {
	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     *int          // Maximum number of retries
}

type HTTPResponseError struct { //nolint:revive // ok with stutter
	Repsonse *http.Response
	Body     []byte
	Cause    error
}

func (e *HTTPResponseError) Error() string {
	code := -1
	if e.Repsonse != nil {
		code = e.Repsonse.StatusCode
	}

	msg := fmt.Sprintf("HTTP Response Error: %d", code)
	if e.Cause != nil {
		msg += fmt.Sprintf(", cause=%v", e.Cause.Error())
	}

	return msg
}

//counterfeiter:generate . Client
type Client interface {
	ConfigureRetries(opts RetryOptions)

	GetJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error)
	PostJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error)
	PutJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error)
	PatchJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error)
	DeleteJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error)
	RequestJSON(ctx context.Context, target interface{}, method, reqURL string, opts ...ClientOpt) (*http.Response, error)

	GetBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error)
	PostBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error)
	PutBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error)
	PatchBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error)
	DeleteBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error)
	RequestBody(ctx context.Context, method, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error)
}

type TelemeterClient struct {
	client        *retryablehttp.Client
	telemeter     *telemetry.Telemeter
	telemeterOpts []telemetry.StartSpanOption
	logger        logging.Logger
}

var (
	_ Client                      = (*TelemeterClient)(nil)
	_ retryablehttp.LeveledLogger = (*HTTPLogger)(nil)
)

func NewTelemeterClient(logger logging.Logger, tel *telemetry.Telemeter, opts ...telemetry.StartSpanOption) *TelemeterClient {
	client := retryablehttp.NewClient()
	if client.HTTPClient.Transport == nil {
		client.HTTPClient.Transport = http.DefaultTransport
	}

	rt := NewTelemeterRoundTripper(client.HTTPClient.Transport, tel, opts...)
	client.HTTPClient.Transport = rt

	client.Logger = &HTTPLogger{
		Logger: logger,
	}

	return &TelemeterClient{
		client:        client,
		telemeter:     tel,
		telemeterOpts: opts,
		logger:        logger,
	}
}

func (c *TelemeterClient) ConfigureRetries(opts RetryOptions) {
	if opts.RetryWaitMin > 0 {
		c.client.RetryWaitMin = opts.RetryWaitMin
	}

	if opts.RetryWaitMax > 0 {
		c.client.RetryWaitMax = opts.RetryWaitMax
	}

	if opts.RetryMax != nil && *opts.RetryMax >= 0 {
		c.client.RetryMax = *opts.RetryMax
	}
}

func (c *TelemeterClient) HTTPClient() *http.Client {
	return c.client.StandardClient()
}

func (c *TelemeterClient) GetJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error) {
	return c.RequestJSON(ctx, target, http.MethodGet, reqURL, opts...)
}

func (c *TelemeterClient) PostJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error) {
	return c.RequestJSON(ctx, target, http.MethodPost, reqURL, opts...)
}

func (c *TelemeterClient) PutJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error) {
	return c.RequestJSON(ctx, target, http.MethodPut, reqURL, opts...)
}

func (c *TelemeterClient) PatchJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error) {
	return c.RequestJSON(ctx, target, http.MethodPatch, reqURL, opts...)
}

func (c *TelemeterClient) DeleteJSON(ctx context.Context, target interface{}, reqURL string, opts ...ClientOpt) (*http.Response, error) {
	return c.RequestJSON(ctx, target, http.MethodDelete, reqURL, opts...)
}

func (c *TelemeterClient) RequestJSON(ctx context.Context, target interface{}, method, reqURL string, opts ...ClientOpt) (*http.Response, error) {
	respHandlerFunc := func(httpResp *http.Response) error {
		if httpResp.StatusCode < 400 {
			err := json.UnmarshalFromReader(httpResp.Body, target)
			return errors.Wrap(err, "could not unmarshal reponse")
		}
		return nil
	}
	opts = append(opts, WithResponseHandler(respHandlerFunc))

	httpResp, err := c.prepareAndSendRequest(ctx, method, reqURL, opts)
	if err != nil {
		return httpResp, errors.Wrap(err, "could not prepareAndSendRequest")
	}
	logger := logging.WithContext(ctx, c.logger)
	logger = logging.WithClientRequest(logger, httpResp.Request)
	defer deferutil.CheckDeferLog(logger, httpResp.Body.Close)

	if httpResp.StatusCode >= 400 {
		body, errBody := io.ReadAll(httpResp.Body)
		if errBody != nil {
			logger.Err("could not read response body", errBody)
		}
		return httpResp, &HTTPResponseError{
			Repsonse: httpResp,
			Body:     body,
		}
	}

	return httpResp, nil
}

func (c *TelemeterClient) GetBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error) {
	return c.RequestBody(ctx, http.MethodGet, reqURL, opts...)
}

func (c *TelemeterClient) PostBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error) {
	return c.RequestBody(ctx, http.MethodPost, reqURL, opts...)
}

func (c *TelemeterClient) PutBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error) {
	return c.RequestBody(ctx, http.MethodPut, reqURL, opts...)
}

func (c *TelemeterClient) PatchBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error) {
	return c.RequestBody(ctx, http.MethodPatch, reqURL, opts...)
}

func (c *TelemeterClient) DeleteBody(ctx context.Context, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error) {
	return c.RequestBody(ctx, http.MethodDelete, reqURL, opts...)
}

func (c *TelemeterClient) RequestBody(ctx context.Context, method, reqURL string, opts ...ClientOpt) ([]byte, *http.Response, error) {
	var body []byte
	respHandlerFunc := func(httpResp *http.Response) error {
		if httpResp.StatusCode < 400 {
			var err error
			body, err = io.ReadAll(httpResp.Body)
			return errors.Wrap(err, "could not read response body")
		}
		return nil
	}
	opts = append(opts, WithResponseHandler(respHandlerFunc))

	httpResp, err := c.prepareAndSendRequest(ctx, method, reqURL, opts)
	if err != nil {
		return nil, httpResp, errors.Wrap(err, "could not prepareAndSendRequest")
	}
	logger := logging.WithContext(ctx, c.logger)
	logger = logging.WithClientRequest(logger, httpResp.Request)
	defer deferutil.CheckDeferLog(logger, httpResp.Body.Close)

	if httpResp.StatusCode >= 400 {
		errBody, errBodyRead := io.ReadAll(httpResp.Body)
		if errBodyRead != nil {
			logger.Err("could not read response body", errBodyRead)
		}
		return nil, httpResp, &HTTPResponseError{
			Repsonse: httpResp,
			Body:     errBody,
		}
	}

	return body, httpResp, nil
}

func (c *TelemeterClient) prepareAndSendRequest(ctx context.Context, method, reqURL string, opts []ClientOpt) (httpResp *http.Response, err error) {
	ctx, span := c.telemeter.StartSpan(ctx, "http", "prepareAndSendRequest", c.telemeterOpts...)
	defer span.End()

	defer func() {
		if httpResp != nil {
			span.SetAttributes(telemetry.HTTPAttributesFromHTTPStatusCode(httpResp)...)
		}

		if err != nil {
			span.SetStatus(telemetry.CodeError, err.Error())
		} else if httpResp != nil {
			code, reason := telemetry.SpanStatusFromHTTPStatusCode(httpResp.StatusCode)
			span.SetStatus(code, reason)
		}
	}()

	req, err := retryablehttp.NewRequestWithContext(ctx, method, reqURL, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	span.SetAttributes(
		telemetry.HTTPClientAttributesFromHTTPRequest(req.Request)...,
	)

	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, errors.Wrap(err, "could not apply request options")
		}
	}

	httpResp, err = c.client.Do(req)
	if err != nil {
		return httpResp, errors.Wrap(err, "could not issue request")
	}

	return httpResp, nil
}
