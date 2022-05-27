package app

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCommand  = "/start"
	finishCommand = "/finish"
)

type App struct {
	bot   *tgbotapi.BotAPI
	users map[int64]worker
}

func NewApp(token string) (App, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return App{}, fmt.Errorf("could not get new telegram bot api: %w", err)
	}

	setCommands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     startCommand,
			Description: "Start pomodoro timer",
		},
		tgbotapi.BotCommand{
			Command:     finishCommand,
			Description: "Finish pomodoro timer",
		},
	)

	if _, err := bot.Request(setCommands); err != nil {
		return App{}, fmt.Errorf("could not set up commands: %w", err)
	}

	return App{
		bot:   bot,
		users: make(map[int64]worker),
	}, nil
}

func (a App) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := a.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			updates.Clear()
			log.Println("Closing Pomodoro bot")
			return
		case update := <-updates:
			if update.Message != nil {
				msg := update.Message
				userID := msg.From.ID
				log.Println("Obtained a message from bot:", msg.Command())

				if msg.Command() == startCommand[1:] {
					if _, ok := a.users[userID]; ok {
						a.SendMessage(msg, "Your Pomodoro timer is already working")
						continue
					}

					newWorker := worker{
						bot:    a.bot,
						chatID: msg.Chat.ID,
						ctx:    make(chan struct{}),
					}
					a.users[userID] = newWorker

					a.SendMessage(msg, "Your Pomodoro timer has started! Now go to work")
					go a.users[userID].Run()

					continue
				}

				if _, ok := a.users[userID]; ok {
					a.SendMessage(msg, "Your Pomodoro timer has stopped. See you next time!")
					a.users[userID].Stop()
					delete(a.users, userID)
					continue
				}

				a.SendMessage(msg, "You don't have a Pomodoro timer launched")
			}
		}
	}
}

func (a App) SendMessage(msg *tgbotapi.Message, msgText string) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	newMsg.ReplyToMessageID = msg.MessageID
	if _, err := a.bot.Send(newMsg); err != nil {
		log.Println("Could not send message:", err)
	}
}
