package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeRate struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type Res struct {
	Bid string `json:"bid"`
}

type ExchangeRateDB struct {
	Code        string    `gorm:"column:code"`
	Codein      string    `gorm:"column:codein"`
	Name        string    `gorm:"column:name"`
	High        float64   `gorm:"column:high"`
	Low         float64   `gorm:"column:low"`
	VarBid      float64   `gorm:"column:varBid"`
	PctChange   float64   `gorm:"column:pctChange"`
	Bid         float64   `gorm:"column:bid"`
	Ask         float64   `gorm:"column:ask"`
	Timestamp   string    `gorm:"column:timestamp"`
	CreatedDate time.Time `gorm:"column:created_date"`
}

func main() {
	http.HandleFunc("/cotacao", HandleGetExchangeRate)
	http.ListenAndServe(":8080", nil)
}

func HandleGetExchangeRate(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	dsn := "root:root@tcp(localhost:3306)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db.AutoMigrate(&ExchangeRateDB{})

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exchangeRateRes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var exchangeRate ExchangeRate
	err = json.Unmarshal(exchangeRateRes, &exchangeRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exchangeRateDB, err := parseToDBFormat(&exchangeRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	db.WithContext(ctx).Create(exchangeRateDB)

	var resJson Res
	resJson.Bid = exchangeRate.USDBRL.Bid

	res, err := json.Marshal(resJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func parseToDBFormat(exchangeRate *ExchangeRate) (*ExchangeRateDB, error) {
	createdDate, err := time.Parse(time.DateTime, exchangeRate.USDBRL.CreateDate)
	if err != nil {
		return nil, err
	}

	high, err := strconv.ParseFloat(exchangeRate.USDBRL.High, 64)
	if err != nil {
		return nil, err
	}

	low, err := strconv.ParseFloat(exchangeRate.USDBRL.Low, 64)
	if err != nil {
		return nil, err
	}

	varBid, err := strconv.ParseFloat(exchangeRate.USDBRL.VarBid, 64)
	if err != nil {
		return nil, err
	}

	pctChange, err := strconv.ParseFloat(exchangeRate.USDBRL.PctChange, 64)
	if err != nil {
		return nil, err
	}

	bid, err := strconv.ParseFloat(exchangeRate.USDBRL.Bid, 64)
	if err != nil {
		return nil, err
	}

	ask, err := strconv.ParseFloat(exchangeRate.USDBRL.Ask, 64)
	if err != nil {
		return nil, err
	}

	return &ExchangeRateDB{
		Code:        exchangeRate.USDBRL.Code,
		Codein:      exchangeRate.USDBRL.Codein,
		Name:        exchangeRate.USDBRL.Name,
		High:        high,
		Low:         low,
		VarBid:      varBid,
		PctChange:   pctChange,
		Bid:         bid,
		Ask:         ask,
		Timestamp:   exchangeRate.USDBRL.Timestamp,
		CreatedDate: createdDate,
	}, nil
}
