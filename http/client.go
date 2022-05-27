package http

import (
	"net/http"

	"github.com/gsmcwhirter/go-util/v10/telemetry"
)

//counterfeiter:generate . Client
type Client interface{}

type TelemeterClient struct {
	client *http.Client
}

func NewTelemeterClient(tel *telemetry.Telemeter, opts ...telemetry.StartSpanOption) *TelemeterClient {
	c := &http.Client{}
	return NewTelemeterClientFrom(c, tel, opts...)
}

func NewTelemeterClientFrom(client *http.Client, tel *telemetry.Telemeter, opts ...telemetry.StartSpanOption) *TelemeterClient {
	rt := NewTelemeterRoundTripper(client.Transport, tel, opts...)
	client.Transport = rt

	return &TelemeterClient{
		client: client,
	}
}

func (c *TelemeterClient) HTTPClient() *http.Client {
	return c.client
}
