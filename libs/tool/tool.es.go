package tool

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"go-schedule/config"

	"github.com/elastic/go-elasticsearch/v8"
)

var fullElasticSearch map[string][]*elasticsearch.Client

func GetElasticSearchClient(key string) *elasticsearch.Client {
	result := fullElasticSearch[key]
	count := GetRandmod(len(result))

	return result[count]
}

func HandleElasticSearchClient() {
	client := make(map[string][]*elasticsearch.Client)

	local := config.GetElasticSearchConfig()

	for k, v := range local {
		m := k + ".master"

		clients := createElasticSearchClient(v.Address, v.Username, v.Password)
		client[m] = append(client[m], clients)
	}

	fullElasticSearch = client
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

func CloseElasticSearch() {
	// for _, val := range fullElasticSearch {
	// 	for _, v := range val {
	// 		// 没找到 close 方法
	// 		v.Close()
	// 	}
	// }

	Stdout("ElasticSearch is Close")
}
