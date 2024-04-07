package http

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/segmentio/kafka-go"
)

func NewAuthRouting(r *fiber.App, a HTTPInterface) {
	r.Post("/v1/ping", a.Get)
}

type HTTPInterface interface {
	Get(ctx fiber.Ctx) error
}

type Handler struct {
	producer	*kafka.Conn
	id 			uint64
	idMutex  *sync.Mutex
	keyChanelMap map[uint64]chan []byte
}

func NewHandler(producer *kafka.Conn, mMap map[uint64]chan []byte) HTTPInterface {
	return &Handler{
		producer: producer,
		id: 0,
		idMutex:  &sync.Mutex{},
		keyChanelMap: mMap,
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
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	outputChanel := make(chan []byte)
	a.keyChanelMap[id] = outputChanel

	responseData := <- outputChanel

	response := &Pong{}
	if err := json.Unmarshal(responseData, response); err != nil {
		fmt.Println(err)
	}

	ctx.JSON(response)
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return ctx.SendStatus(fiber.StatusOK)
}

func Reponser(
	ctx context.Context,
	concumer *kafka.Reader,
	keyChanelMap map[uint64]chan []byte,
) {
	for {
		message, err := concumer.ReadMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		}
		
		num := binary.LittleEndian.Uint64(message.Key)
		outputChanel := keyChanelMap[num]

		outputChanel <- message.Value
	}

	if err := concumer.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}