package tool

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Fetch struct {
	client *http.Client
}

func NewFetch() *Fetch {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	transport.TLSHandshakeTimeout = 3 * time.Second

	return &Fetch{
		client: &http.Client{
			Timeout:   4 * time.Second,
			Transport: transport,
		},
	}
}

func (f *Fetch) GET(uri string, params map[string]string) ([]byte, error) {
	return f.doRequest(http.MethodGet, uri, nil, params, nil)
}

func (f *Fetch) POST(uri string, body any, params map[string]string) ([]byte, error) {
	var buffer io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buffer = bytes.NewReader(data)
	}

	return f.doRequest(http.MethodPost, uri, buffer, params, nil)
}

func (f *Fetch) doRequest(method, uri string, body io.Reader, params, headers map[string]string) ([]byte, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if params != nil {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	// record := fmt.Sprintf("url:%s, result:%s", req.URL.String(), string(data))
	// fmt.Println("->", record, err)

	return data, err
}
