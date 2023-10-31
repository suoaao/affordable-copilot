package httputil

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(targetUrl string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return nil, err
	}
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}
	return &httputil.ReverseProxy{
		Director:  director,
		Transport: &CachedTransport{http.DefaultTransport},
	}, nil
}

type CachedTransport struct {
	http.RoundTripper
}

func (t *CachedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.RoundTripper.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return response, nil
	}
	portal, ok := request.Context().Value("is_cache").(chan []byte)
	if !ok {
		return response, nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body = io.NopCloser(bytes.NewReader(body))

	// cache response body
	portal <- body

	return response, nil
}
