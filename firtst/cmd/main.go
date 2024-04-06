package main

import (
	"context"
	"fmt"
	"log"

	"first/http"
	"first/pkg/conf"
	"first/pkg/kafka"

	"github.com/gofiber/fiber/v3"
)

func main() {
	config, err := conf.ReadConf("./pkg/conf/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	producer, concumer, err := kafka.KafkaConn(config)
	if err != nil {
		log.Fatal(err)
	}

	keyChanelMap := make(map[uint64]chan []byte)

	go http.Reponser(context.Background(), concumer, keyChanelMap)

	router := fiber.New()
	http.NewAuthRouting(router, http.NewHandler(producer, keyChanelMap))

	log.Fatal(router.Listen(fmt.Sprintf(":%d", config.Main.HTTPPort)))
}