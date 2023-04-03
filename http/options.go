package http

import (
	"io"
	"net/http"
	"net/url"
)

type ClientOpt = func(req *Request) error

func WithBody(body io.ReadCloser) ClientOpt {
	return func(req *Request) error {
		return req.SetBody(body)
	}
}

func WithHeaders(headers http.Header) ClientOpt {
	return func(req *Request) error {
		for k, vs := range headers {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}

		return nil
	}
}

func WithQueryParams(params url.Values) ClientOpt {
	return func(req *Request) error {
		query := req.URL.Query()
		for k, vs := range params {
			for _, v := range vs {
				query.Add(k, v)
			}
		}
		req.URL.RawQuery = query.Encode()

		return nil
	}
}

func WithResponseHandler(fn ResponseHandlerFunc) ClientOpt {
	return func(req *Request) error {
		req.SetResponseHandler(fn)
		return nil
	}
}
