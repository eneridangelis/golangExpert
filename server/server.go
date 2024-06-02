package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	Code        string
	Codein      string
	Name        string
	High        float64
	Low         float64
	VarBid      float64
	PctChange   float64
	Bid         float64
	Ask         float64
	Timestamp   string
	CreatedDate time.Time
}

func main() {
	http.HandleFunc("/cotacao", HandleGetExchangeRate)
	http.ListenAndServe(":8080", nil)
}

func HandleGetExchangeRate(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

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

	err = InsertExchangeRate(ctx, db, exchangeRateDB)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resJson Res
	resJson.Bid = exchangeRate.USDBRL.Bid

	res, err := json.Marshal(resJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func InsertExchangeRate(ctx context.Context, db *sql.DB, exchangeRate *ExchangeRateDB) error {
	stmt, err := db.PrepareContext(ctx, "INSERT INTO exchange_rate (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, created_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, exchangeRate.Code, exchangeRate.Codein, exchangeRate.Name, exchangeRate.High, exchangeRate.Low, exchangeRate.VarBid, exchangeRate.PctChange, exchangeRate.Bid, exchangeRate.Ask, exchangeRate.Timestamp, exchangeRate.CreatedDate)
	if err != nil {
		return err
	}

	return nil
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
