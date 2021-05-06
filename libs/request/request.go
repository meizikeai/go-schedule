package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"go-schedule/libs/tool"
	"go-schedule/libs/types"

	log "github.com/sirupsen/logrus"
)

func GET(reqUrl string, reqParams types.MapStringString, headers types.MapStringString) ([]byte, error) {
	var req *http.Request
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

	req, err = http.NewRequest("GET", url, nil)

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

	pool := &http.Client{}
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

	record := types.MapStringInterface{
		"url":    urlPath.String(),
		"result": result,
	}

	log.Info(string(tool.MarshalJson(record)))

	return result, err
}

func POST(reqUrl string, body types.MapStringInterface, params types.MapStringString, headers types.MapStringString) ([]byte, error) {
	var req *http.Request
	result := []byte{}

	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(data))

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

	pool := &http.Client{}
	res, err := pool.Do(req)

	if err != nil {
		log.Error(err)
		return result, err
	}

	defer res.Body.Close()

	result, err = ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error(err)
		return result, err
	}

	record := types.MapStringInterface{
		"url":    req.URL.String(),
		"body":   data,
		"result": result,
	}

	log.Info(string(tool.MarshalJson(record)))

	return result, err
}
