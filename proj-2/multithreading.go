package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	c1 := make(chan string)
	c2 := make(chan string)

	req1, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/01153000+cep", nil)
	if err != nil {
		panic(err)
	}

	req2, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", "01153000"), nil)
	if err != nil {
		panic(err)
	}

	go func() {
		resp, err := http.DefaultClient.Do(req1)
		if err != nil {
			panic(err)
		}

		cep, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		c1 <- string(cep)
	}()

	go func() {
		resp, err := http.DefaultClient.Do(req2)
		if err != nil {
			panic(err)
		}

		cep, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		c2 <- string(cep)
	}()

	select {
	case msg1 := <-c1:
		fmt.Println("Received from Brasil API first:", msg1)
	case msg2 := <-c2:
		fmt.Println("Received from ViaCEP first:", msg2)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}

}
