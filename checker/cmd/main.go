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

type Ping struct {
	Ping string `json:"ping"`
}

type Pong struct {
	Ping string `json:"ping"`
	Pong string `json:"pong"`
}

func check(
	ctx context.Context,
	client *resty.Client,
	t time.Time,
) {
	val := strconv.Itoa(rand.Int())
		requestData := &Ping{Ping: val}
	
		request, err := json.Marshal(requestData)
		if err != nil {
			fmt.Println(time.Since(t), err)
			return
		}
		
		resp, err := client.R().
			SetContext(ctx).
			SetBody(request).
			Post("http://127.0.0.1:8000/v1/ping")
		if err != nil {
			fmt.Println(time.Since(t), err)
			return
		}

		if resp.StatusCode() != 200 {
			fmt.Println("status not 200")
			return
		}
	
		response := &Pong{}
	
		if err = json.Unmarshal(resp.Body(), response); err != nil {
			fmt.Println(time.Since(t), err)
			return
		}
	
		if response.Ping != val || response.Pong != val {
			fmt.Println(time.Since(t), false)
			return
		}
}


func checker(
	ctx context.Context,
	client *resty.Client,
	t time.Time,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			check(ctx, client, t)
		}
	}
}

func main() {
	client := resty.New()

	t := time.Now()

	wg := &sync.WaitGroup{}

	ctx, cancel :=  context.WithTimeout(context.Background(), 10 * time.Minute)
	defer cancel()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go checker(ctx, client, t, wg)
	}

	wg.Wait()
}