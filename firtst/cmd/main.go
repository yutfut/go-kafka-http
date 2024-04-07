package main

import (
	"context"
	"fmt"
	"log"

	"first/http"
	"first/pkg/conf"
	"first/pkg/kafka"
	"first/pkg/syncmap"

	"github.com/gofiber/fiber/v3"
)

func main() {
	config, err := conf.ReadConf("./pkg/conf/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	producer, consumer, err := kafka.KafkaConn(config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}	
	} ()

	// keyChanelMap := make(map[uint64]chan []byte)

	sMap := syncmap.NewSyncMap()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go http.Responser(ctx, consumer, sMap)

	router := fiber.New()
	http.NewRouting(router, http.NewHandler(producer, sMap))

	log.Fatal(router.Listen(fmt.Sprintf(":%d", config.Main.HTTPPort)))
}