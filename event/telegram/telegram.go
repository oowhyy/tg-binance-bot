package telegram

import (
	"errors"
	"fmt"

	"github.com/oowhyy/tg-binance-bot/calc"
	"github.com/oowhyy/tg-binance-bot/client/telegram"
	"github.com/oowhyy/tg-binance-bot/event"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrWrongMetaType    = errors.New("wrong meta type")
)

type EventMeta struct {
	ChatId   int
	Username string
}

// implements both Interfaces - Fetcher and Processor
type Operator struct {
	tg         *telegram.Client
	calculator *calc.Calc
	offset     int
}

func New(tgCl *telegram.Client, calc *calc.Calc) *Operator {
	return &Operator{
		tg:         tgCl,
		calculator: calc,
	}
}

// implementation of Fetcher interface
func (op *Operator) Fetch(limit int) ([]event.Event, error) {
	updates, err := op.tg.Update(op.offset, limit) // get client updates
	if err != nil {
		return nil, fmt.Errorf("unable to fetch events: %w", err)
	}
	n := len(updates)
	if n == 0 {
		return nil, nil
	}
	res := make([]event.Event, len(updates))
	for k, upd := range updates {
		res[k] = toEvent(upd)
	}
	op.offset = updates[n-1].Id + 1
	return res, nil
}

// implementation of Processor interface
func (op *Operator) Process(ev event.Event) error {
	switch ev.Type {
	case event.Message:
		return op.processMessage(ev)
	default:
		return ErrUnknownEventType
	}
}

// processes event of Message type
func (op *Operator) processMessage(event event.Event) error {
	meta, ok := event.Meta.(EventMeta)
	if !ok {
		return ErrWrongMetaType
	}
	if err := op.doCommand(event.Text, meta.ChatId, meta.Username); err != nil {
		return fmt.Errorf("unable to do command: %w", err)
	}
	return nil
}

// converts update to event
func toEvent(upd telegram.Update) event.Event {
	updType := eventType(upd)
	res := event.Event{
		Type: updType,
		Text: eventText(upd),
	}
	if updType == event.Message {
		res.Meta = EventMeta{
			ChatId:   upd.Message.Chat.Id,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

// gets event text
func eventText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

// gets event type
func eventType(upd telegram.Update) event.Type {
	if upd.Message == nil {
		return event.Unknown
	}
	return event.Message
}
