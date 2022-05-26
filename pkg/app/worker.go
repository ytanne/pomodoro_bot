package app

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type condition int

const (
	working = iota
	chilling

	workingPeriod  = 25 * time.Minute
	chillingPeriod = 5 * time.Minute

	workingMsg  = "It's working time"
	chillingMsg = "It's time to chill"
)

type worker struct {
	current condition
	bot     *tgbotapi.BotAPI
	chatID  int64
	ctx     chan struct{}
}

func (w worker) Run() {
	w.current = working
	ticker := time.NewTicker(workingPeriod)

	for {
		select {
		case <-w.ctx:
			ticker.Stop()
			return
		case <-ticker.C:
			if w.current == condition(working) {
				w.current = chilling

				w.SendMessage(chillingMsg)
				ticker.Reset(chillingPeriod)
				continue
			}

			w.current = working
			w.SendMessage(workingMsg)
			ticker.Reset(workingPeriod)
		}
	}
}

func (w worker) SendMessage(msg string) {
	newMsg := tgbotapi.NewMessage(w.chatID, msg)

	if _, err := w.bot.Send(newMsg); err != nil {
		log.Println("could not send message:", err)
	}
}

func (w worker) Stop() {
	w.ctx <- struct{}{}
}
