package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-schedule/libs/tool"
	"go-schedule/libs/types"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func Get(reqUrl string, reqParams types.MapStringString) ([]byte, error) {
	params := url.Values{}
	urlPath, err := url.Parse(reqUrl)

	if err != nil {
		log.Error(err)
		return []byte{}, err
	}

	for key, val := range reqParams {
		params.Set(key, val)
	}

	urlPath.RawQuery = params.Encode()
	res, err := http.Get(urlPath.String())

	if err != nil {
		log.Error(err)
		return []byte{}, err
	}

	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error(err)
	}

	record := types.MapStringInterface{
		"url":    urlPath.String(),
		"result": result,
	}

	log.Info(string(tool.MarshalJson(record)))

	return result, err
}

func Post(url string, body types.MapStringInterface, params types.MapStringString, headers types.MapStringString) ([]byte, error) {
	var req *http.Request

	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

	if err != nil {
		log.Error(err)
		return nil, errors.New("new request is fail")
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
		return []byte{}, err
	}

	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Error(err)
	}

	record := types.MapStringInterface{
		"url":    req.URL.String(),
		"body":   data,
		"result": result,
	}

	log.Info(string(tool.MarshalJson(record)))

	return result, err
}
