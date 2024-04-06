package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Ping struct {
	Ping string `json:"ping"`
}

type Pong struct {
	Ping string `json:"ping"`
	Pong string `json:"pong"`
}

func main() {
	ctx := context.Background()
	producer, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", "output", 0)
	if err != nil {
		log.Fatal(err)
	}

	concumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		GroupID:   "consumer-group-id",
		Topic:     "input",
		MaxBytes:  10e6, // 10MB
	})

	for {
		message, err := concumer.ReadMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		}

		request := &Ping{}
		if err = json.Unmarshal(message.Value, request); err != nil {
			fmt.Println(err)
		}

		fmt.Println(request)

		response := &Pong{
			Ping: request.Ping,
			Pong: request.Ping,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(jsonResponse)

		if _, err = producer.WriteMessages(kafka.Message{Key: message.Key, Value: jsonResponse}); err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}

	if err := concumer.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}