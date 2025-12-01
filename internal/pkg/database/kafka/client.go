// internal/pkg/database/kafka/kafka.go
package kafka

import (
	"context"
	"fmt"
	"time"

	"go-schedule/internal/config"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	writers map[string]*kafka.Writer
	readers map[string]*kafka.Reader
}

func NewClient(cfgs *map[string]config.KafkaInstance) *Client {
	c := &Client{
		writers: make(map[string]*kafka.Writer),
		readers: make(map[string]*kafka.Reader),
	}

	for name, cfg := range *cfgs {
		w := &kafka.Writer{
			Addr:                   kafka.TCP(cfg.Brokers...),
			Topic:                  cfg.Topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
			BatchSize:              100,
			BatchTimeout:           10 * time.Millisecond,
			ReadTimeout:            10 * time.Second,
			WriteTimeout:           10 * time.Second,
			RequiredAcks:           kafka.RequireOne,
			Async:                  false,
			// Transport: &kafka.Transport{
			// 	TLS: &tls.Config{},
			// },
		}
		c.writers[name] = w

		if cfg.GroupID != "" {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:  cfg.Brokers,
				GroupID:  cfg.GroupID,
				Topic:    cfg.Topic,
				MinBytes: 10e3, // 10KB
				MaxBytes: 10e6, // 10MB
			})
			c.readers[name] = r
		}
	}
	return c
}

func (c *Client) Produce(ctx context.Context, name string, key, value []byte) error {
	w, ok := c.writers[name]
	if !ok {
		return fmt.Errorf("kafka writer %s not found", name)
	}
	return w.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

func (c *Client) Consume(ctx context.Context, name string) (kafka.Message, error) {
	r, ok := c.readers[name]
	if !ok {
		return kafka.Message{}, fmt.Errorf("kafka reader %s not found", name)
	}
	return r.FetchMessage(ctx)
}

func (c *Client) Commit(ctx context.Context, name string, msg kafka.Message) error {
	r, ok := c.readers[name]
	if !ok {
		return nil
	}
	return r.CommitMessages(ctx, msg)
}

func (c *Client) Close() {
	for _, w := range c.writers {
		w.Close()
	}
	for _, r := range c.readers {
		r.Close()
	}
}
