package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Rate struct {
	Base      string
	Timestamp time.Time
	Quotes    map[string]float32
}

var (
	rateChannel             = make(chan []*Rate)
	ExchangeServiceProvider ExchangeRateService
)

type ExchangeRateService interface {
	GetLatestRates(ctx context.Context) ([]*Rate, error)
}

func main() {
	ctx := context.Background()
	var currentRates []*Rate

	go func() {
		select {
		case val := <-rateChannel:
			currentRates = val
		case <-time.After(3 * time.Minute):
			fmt.Println("timed out")
		}
	}()

	for {
		rates, err := ExchangeServiceProvider.GetLatestRates(ctx)
		if err != nil {
			log.Printf("Unable to get latest rates: %v \n", err)
			continue
		}
		rateChannel <- rates
		time.Sleep(time.Minute)
	}
}
