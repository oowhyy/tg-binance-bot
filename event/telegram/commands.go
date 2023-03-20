package telegram

import (
	"fmt"
	"log"
	"strings"

	"github.com/oowhyy/tg-binance-bot/calc"
)

const (
	startCmd = "/start"
	nameCmd  = "/username"
	helpCmd  = "/help"
	bnCmd    = "/triangle"
)

func (op *Operator) doCommand(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from %s", text, username)
	switch text {
	case startCmd:
		return op.tg.Message(chatId, msgHello+msgHelp)
	case nameCmd:
		return op.tg.Message(chatId, "Hi, "+username)
	case bnCmd:
		return op.triangles(chatId)
	case helpCmd:
		return op.tg.Message(chatId, msgHelp)
	default:
		return op.tg.Message(chatId, msgUnknownCommand)
	}
}

func (op *Operator) triangles(chatId int) error {

	res, err := op.calculator.BestTriangles(calc.DefaultTrianglesNum)
	if err != nil {
		return err
	}
	formated := trianglesToText(res)
	return op.tg.Message(chatId, formated)

}

func trianglesToText(res []calc.Triangle) string {
	sb := &strings.Builder{}
	for _, t := range res {
		for _, sy := range t.Coins {
			sb.WriteString(fmt.Sprintf("%s -> ", sy))
		}
		// sb.WriteString(fmt.Sprintf("%s", t.Coins[0]))
		sb.WriteString(t.Coins[0])
		sb.WriteString(fmt.Sprintf("  %.6f%%\n", 100*t.Coef))
	}
	return sb.String()
}
