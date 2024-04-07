package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// func checker(client *resty.Client) {

// }

type Ping struct {
	Ping string `json:"ping"`
}

type Pong struct {
	Ping string `json:"ping"`
	Pong string `json:"pong"`
}

func checker(
	ctx context.Context,
	client *resty.Client,
	t time.Time,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		val := strconv.Itoa(rand.Int())
		reqestData := &Ping{Ping: val}
	
		reqest, err := json.Marshal(reqestData)
		if err != nil {
			fmt.Println(time.Since(t), err)
		}
		
		resp, err := client.R().
			SetContext(ctx).
			SetBody(reqest).
			Post("http://127.0.0.1:8000/v1/ping")
		if err != nil {
			fmt.Println(time.Since(t), err)
		}
	
		response := &Pong{}
	
		if err = json.Unmarshal(resp.Body(), response); err != nil {
			fmt.Println(time.Since(t), err)
		}
	
		if response.Ping != val || response.Pong != val {
			fmt.Println(time.Since(t), false)
		}
	}
}

func main() {
	client := resty.New()

	t := time.Now()

	wg := &sync.WaitGroup{}

	ctx, cancel :=  context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go checker(ctx, client, t, wg)
	}

	wg.Wait()
}