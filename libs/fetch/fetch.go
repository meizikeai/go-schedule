package fetch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go-schedule/libs/types"

	log "github.com/sirupsen/logrus"
)

func GET(reqUrl string, reqParams, headers types.MapStringString) ([]byte, error) {
	result := []byte{}

	params := url.Values{}
	urlPath, err := url.Parse(reqUrl)

	if err != nil {
		log.Error(err)
		return result, err
	}

	for key, val := range reqParams {
		params.Set(key, val)
	}

	urlPath.RawQuery = params.Encode()
	url := urlPath.String()

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Error(err)
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	pool := &http.Client{
		Timeout: 200 * time.Millisecond,
	}

	res, err := pool.Do(req)

	if err != nil {
		log.Error(err)
		return result, err
	}

	result, err = ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error(err)
		return result, err
	}

	record := fmt.Sprintf("url:%s, result:%s", req.URL.String(), string(result))
	log.Info(record)

	defer res.Body.Close()

	return result, err
}

func POST(reqUrl string, body interface{}, params, headers types.MapStringString) ([]byte, error) {
	result := []byte{}

	data, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(data))

	if err != nil {
		log.Error(err)
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
		Timeout: 200 * time.Millisecond,
	}

	res, err := pool.Do(req)

	if err != nil {
		log.Error(err)
		return result, err
	}

	result, err = ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error(err)
		return result, err
	}

	record := fmt.Sprintf("url:%s, body:%s, result:%s", req.URL.String(), string(data), string(result))
	log.Info(record)

	defer res.Body.Close()

	return result, err
}

func DELETE(reqUrl string, body interface{}, params, headers types.MapStringString) ([]byte, error) {
	result := []byte{}

	data, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(data))

	if err != nil {
		log.Error(err)
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
		Timeout: 200 * time.Millisecond,
	}

	res, err := pool.Do(req)

	if err != nil {
		log.Error(err)
		return result, err
	}

	result, err = ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error(err)
		return result, err
	}

	record := fmt.Sprintf("url:%s, body:%s, result:%s", req.URL.String(), string(data), string(result))
	log.Info(record)

	defer res.Body.Close()

	return result, err
}
