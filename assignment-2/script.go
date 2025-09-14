package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/linkedin/goavro/v2"
	"github.com/segmentio/kafka-go"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <topic> <consumer_group>")
	}

	topic := os.Args[1]
	consumerGroup := os.Args[2]

	// Load CA certificate
	caCert, err := os.ReadFile("auth/ca.crt")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to parse CA certificate")
	}

	// Load client certificate and key
	cert, err := tls.LoadX509KeyPair("auth/kafka-cert.pem", "auth/kafka-key.pem")
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

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
			TLS: tlsConfig,
		},
	})

	defer r.Close()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		fmt.Println("\nEXITING SAFELY!")
		r.Close()
		os.Exit(0)
	}()

	fmt.Printf("Consumer started for topic: %s, group: %s\n", topic, consumerGroup)

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		// Decode OCF message
		ocfReader, err := goavro.NewOCFReader(bytes.NewReader(msg.Value))
		if err != nil {
			log.Printf("Failed to create OCF reader: %v", err)
			fmt.Printf("Raw message: %x\n", msg.Value)
			continue
		}

		// Read all records from the OCF container
		for ocfReader.Scan() {
			record, err := ocfReader.Read()
			if err != nil {
				log.Printf("Failed to read record: %v", err)
				break
			}

			// Print the schema name
			schema := ocfReader.Codec().Schema()
			var schemaMap map[string]interface{}
			if err := json.Unmarshal([]byte(schema), &schemaMap); err != nil {
				log.Printf("Failed to parse schema: %v", err)
			} else {
				fmt.Printf("%s\n", schemaMap["name"].(string))
			}

			jsonData, err := json.Marshal(record)
			if err != nil {
				log.Printf("Failed to marshal record: %v", err)
				continue
			}
			fmt.Printf("%s\n", string(jsonData))
		}

		if err := ocfReader.Err(); err != nil {
			log.Printf("OCF reader error: %v", err)
		}
	}
}
