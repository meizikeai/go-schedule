package tool

import (
	"context"
	"fmt"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

type EtchType struct {
	Client *etcd.Client
}

func NewEtcdServer(address []string, username, password string) *EtchType {
	config := etcd.Config{
		Endpoints:   address,
		Username:    username,
		Password:    password,
		DialTimeout: 4 * time.Second,
	}

	client, err := etcd.New(config)

	if err != nil {
		panic(err)
	}

	return &EtchType{
		Client: client,
	}
}

func (e *EtchType) Put(key, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	_, err := e.Client.Put(ctx, key, value)
	cancel()

	if err != nil {
		fmt.Printf("put to etcd failed, err:%v", err)
	}
}

func (e *EtchType) Get(key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	res, err := e.Client.Get(ctx, key)
	cancel()

	result := ""

	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return result
	}

	for _, ev := range res.Kvs {
		result = string(ev.Value)
	}

	return result
}
