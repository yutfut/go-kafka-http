package http

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"first/pkg/syncmap"

	"github.com/gofiber/fiber/v3"
	"github.com/segmentio/kafka-go"
)

func NewRouting(r *fiber.App, a HTTPInterface) {
	r.Post("/v1/ping", a.Get)
}

type HTTPInterface interface {
	Get(ctx fiber.Ctx) error
}

type Handler struct {
	producer	*kafka.Conn
	id 			uint64
	idMutex  *sync.Mutex
	keyChanelMap syncmap.SyncMapInterface
}

func NewHandler(producer *kafka.Conn, sMap syncmap.SyncMapInterface) HTTPInterface {
	return &Handler{
		producer: producer,
		id: 0,
		idMutex:  &sync.Mutex{},
		keyChanelMap: sMap,
	}
}

type Ping struct {
	Ping string `json:"ping"`
}

type Pong struct {
	Ping string `json:"ping"`
	Pong string `json:"pong"`
}

func (a *Handler) getID() uint64 {
	a.idMutex.Lock()
	ID := a.id
	a.id += 1
	a.idMutex.Unlock()
	return ID
}

func (a *Handler)Get(ctx fiber.Ctx) error {
	id := a.getID()

	byteID := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteID, id)

	if _, err := a.producer.WriteMessages(kafka.Message{Key: byteID, Value: ctx.Body()}); err != nil {
		fmt.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	outputChanel := make(chan []byte)
	if err := a.keyChanelMap.Set(id, outputChanel); err != nil {
		fmt.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	responseData := <- outputChanel

	response := &Pong{}
	if err := json.Unmarshal(responseData, response); err != nil {
		fmt.Println(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	ctx.JSON(response)
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return ctx.SendStatus(fiber.StatusOK)
}

func Responser(
	ctx context.Context,
	consumer *kafka.Reader,
	keyChanelMap syncmap.SyncMapInterface,
) {
	for {
		message, err := consumer.ReadMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		}
		
		num := binary.LittleEndian.Uint64(message.Key)
		
		outputChanel, err := keyChanelMap.Get(num)
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func () {
			outputChanel <- message.Value
			if err = keyChanelMap.Del(num); err != nil {
				fmt.Println(err)
			}
		} ()
	}

	if err := consumer.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}