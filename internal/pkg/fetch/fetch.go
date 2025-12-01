package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Fetch struct {
	client      *http.Client
	baseHeaders map[string]string
}

func NewFetch() *Fetch {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	transport.TLSHandshakeTimeout = 5 * time.Second
	transport.MaxIdleConns = 100
	transport.MaxIdleConnsPerHost = 20
	transport.IdleConnTimeout = 90 * time.Second

	// OpenTelemetry
	// tr := otelhttp.NewTransport(transport)

	return &Fetch{
		client: &http.Client{
			Timeout:   4 * time.Second,
			Transport: transport,
		},
		baseHeaders: map[string]string{
			"Content-Type": "application/json",
			// "User-Agent":   "go-practice/1.0",
		},
	}
}

func (f *Fetch) GET(ctx context.Context, uri string, params map[string]string) ([]byte, int, error) {
	return f.do(ctx, http.MethodGet, uri, nil, params, nil)
}

func (f *Fetch) POST(ctx context.Context, uri string, body any) ([]byte, int, error) {
	var reader io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reader = bytes.NewReader(data)
	}
	return f.do(ctx, http.MethodPost, uri, reader, nil, nil)
}

func (f *Fetch) do(ctx context.Context, method, rawURL string, body io.Reader, params, headers map[string]string) ([]byte, int, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, 0, err
	}
	if params != nil {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, 0, err
	}

	for k, v := range f.baseHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// if gctx, ok := ginctx.FromContext(ctx); ok {
	// 	if reqID := gctx.GetReqID(); reqID != "" {
	// 		req.Header.Set("X-Request-Id", reqID)
	// 	}
	// }

	resp, err := f.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	// record := fmt.Sprintf("url:%s, result:%s", req.URL.String(), string(data))
	// fmt.Println("->", record)

	return data, resp.StatusCode, nil
}
