package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "https://httpbin.org/get", "URL do serviço a ser testado")
	totalRequests := flag.Int("requests", 100, "Número total de requests")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		log.Fatal("URL não fornecida")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var successfulRequests, totalRequestsMade int
	statusCounts := make(map[int]int)

	startTime := time.Now()

	sem := make(chan struct{}, *concurrency)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	for i := 0; i < *totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			resp, err := client.Get(*url)
			mu.Lock()
			if err != nil {
				log.Printf("Erro ao fazer request: %v", err)
			} else {
				totalRequestsMade++
				if resp.StatusCode == http.StatusOK {
					successfulRequests++
				}
				statusCounts[resp.StatusCode]++
			}
			mu.Unlock()
		}()
	}

	select {
	case <-done:
	case <-time.After(10 * time.Minute):
		log.Println("Timeout global alcançado")
	}

	duration := time.Since(startTime)

	fmt.Printf("Tempo total gasto: %v\n", duration)
	fmt.Printf("Quantidade total de requests realizados: %d\n", totalRequestsMade)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", successfulRequests)
	fmt.Printf("Distribuição de outros códigos de status HTTP:\n")
	for status, count := range statusCounts {
		fmt.Printf("  %d: %d\n", status, count)
	}
}
