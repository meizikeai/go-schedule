package fetch

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

func GET(reqUrl string, reqParams, headers map[string]string) ([]byte, error) {
	result := []byte{}

	params := url.Values{}
	urlPath, err := url.Parse(reqUrl)

	if err != nil {
		return result, err
	}

	for key, val := range reqParams {
		params.Set(key, val)
	}

	urlPath.RawQuery = params.Encode()
	url := urlPath.String()

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	pool := &http.Client{
		Timeout: 4000 * time.Millisecond,
	}

	res, err := pool.Do(req)

	if err != nil {
		return result, err
	}

	result, err = io.ReadAll(res.Body)

	if err != nil {
		return result, err
	}

	// record := fmt.Sprintf("url:%s, result:%s", req.URL.String(), string(result))
	// fmt.Println(record)

	defer res.Body.Close()

	return result, nil
}

func POST(reqUrl string, body any, params, headers map[string]string) ([]byte, error) {
	result := []byte{}

	data, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(data))

	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()

	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}

		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	pool := &http.Client{
		Timeout: 4000 * time.Millisecond,
	}

	res, err := pool.Do(req)

	if err != nil {
		return result, err
	}

	result, err = io.ReadAll(res.Body)

	if err != nil {
		return result, err
	}

	// record := fmt.Sprintf("url:%s, body:%s, result:%s", req.URL.String(), string(data), string(result))
	// fmt.Println(record)

	defer res.Body.Close()

	return result, nil
}
