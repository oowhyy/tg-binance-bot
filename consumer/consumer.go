package consumer

import (
	"log"
	"time"

	"github.com/oowhyy/tg-binance-bot/event"
)

type Consumer struct {
	fetcher   event.Fetcher
	processor event.Processor
	batchSize int
}

func New(fetcher event.Fetcher, processor event.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher,
		processor,
		batchSize,
	}
}

func (cons *Consumer) Start() error {
	for {
		events, err := cons.fetcher.Fetch(cons.batchSize)
		if err != nil {
			log.Printf("consumer fetch error %s", err.Error())
			continue
		}
		if len(events) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := cons.handleEvents(events); err != nil {
			log.Print(err)
			return err
		}

	}
}

func (cons *Consumer) handleEvents(events []event.Event) error {
	for _, ev := range events {
		log.Printf("got new event: %s", ev.Text)
		if err := cons.processor.Process(ev); err != nil {
			log.Printf("unable to process event %s", err.Error())
		}
	}
	return nil
}
