package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "client36"
	consumerGroup := "my-consumer-group"

	// Create a new reader with SSL configuration
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			"kafka1.dlandau.nl:19092",
			"kafka2.dlandau.nl:29092",
			"kafka3.dlandau.nl:39092",
		},
		Topic:       topic,
		GroupID:     consumerGroup,
		StartOffset: kafka.LastOffset,
		Dialer: &kafka.Dialer{
			TLS: &tls.Config{
				InsecureSkipVerify: true, // equivalent to ssl.endpoint.identification.algorithm: 'none'
			},
		},
	})

	defer r.Close()

	fmt.Printf("Consumer started for topic: %s, group: %s\n", topic, consumerGroup)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}
		fmt.Printf("%s\n", string(msg.Value))
	}
}
