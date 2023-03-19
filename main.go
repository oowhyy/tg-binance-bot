package main

import (
	"flag"
	"log"

	"github.com/oowhyy/tg-binance-bot/calc"
	"github.com/oowhyy/tg-binance-bot/client/binance"
	tgCl "github.com/oowhyy/tg-binance-bot/client/telegram"
	"github.com/oowhyy/tg-binance-bot/consumer"
	tgEv "github.com/oowhyy/tg-binance-bot/event/telegram"
	// "github.com/oowhyy/tg-binance-bot/event"
)

const (
	tgHost    = "api.telegram.org"
	bnHost    = "api.binance.com"
	batchSize = 10
)

func main() {
	calc := calc.New(binance.New(bnHost))
	eventOperator := tgEv.New(tgCl.New(tgHost, mustToken()), calc)
	log.Println("service started")
	consumer := consumer.New(eventOperator, eventOperator, batchSize)
	// bnClient := binance.New(bnHost)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

func mustToken() string {
	token := flag.String("t", "", "token is used to access tg bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("empty token")
	}
	return *token
}
