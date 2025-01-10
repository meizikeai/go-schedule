package tool

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"go-schedule/libs/types"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearch struct {
	Client map[string][]*elasticsearch.Client
}

func NewElasticSearch(data map[string]types.ConfElasticSearch) *ElasticSearch {
	client := make(map[string][]*elasticsearch.Client)

	for k, v := range data {
		m := k + ".master"

		clients := createElasticSearchClient(v.Address, v.Username, v.Password)
		client[m] = append(client[m], clients)
	}

	return &ElasticSearch{
		Client: client,
	}
}

func createElasticSearchClient(address []string, username, password string) *elasticsearch.Client {
	config := elasticsearch.Config{
		Addresses: address,
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 4 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 4 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	client, err := elasticsearch.NewClient(config)

	if err != nil {
		panic(err)
	}

	_, err = client.Info()

	if err != nil {
		panic(err)
	}

	return client
}

func (e *ElasticSearch) Close() {
	// for _, val := range e.Client {
	// 	for _, v := range val {
	// 		// Can't find a way to close it
	// 		v.Close()
	// 	}
	// }
}
