package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Res struct {
	Bid string `json:"bid"`
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	exchangeRate, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var res Res
	err = json.Unmarshal(exchangeRate, &res)
	if err != nil {
		return
	}

	log.Printf("bid: %s", res.Bid)

	f, err := os.Create("cotacao.txt")
	if err != nil {
		return
	}
	defer f.Close()

	value := fmt.Sprintf("Dólar:%s", res.Bid)
	f.WriteString(value)
}
