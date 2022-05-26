package app

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ytanne/pomodoro_bot/pkg/config"
)

const (
	startCommand  = "/start"
	finishCommand = "/finish"
)

type App struct {
	bot   *tgbotapi.BotAPI
	users map[string]worker
}

func NewApp(cfg config.Config) (App, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return App{}, fmt.Errorf("could not get new telegram bot api: %w", err)
	}

	setCommands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     startCommand,
			Description: "Starting pomodoro timer",
		},
		tgbotapi.BotCommand{
			Command:     finishCommand,
			Description: "Finishing pomodoro timer",
		},
	)

	if _, err := bot.Request(setCommands); err != nil {
		return App{}, fmt.Errorf("could not set up commands: %w", err)
	}

	return App{
		bot:   bot,
		users: make(map[string]worker),
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
				userName := msg.From.UserName
				log.Println("Obtained a message from bot:", msg.Command())

				if msg.Command() == startCommand[1:] {
					if _, ok := a.users[userName]; ok {
						a.SendMessage(msg, "Your Pomodoro timer is already working")
						continue
					}

					newWorker := worker{
						bot:    a.bot,
						chatID: msg.Chat.ID,
						ctx:    make(chan struct{}),
					}
					a.users[userName] = newWorker

					a.SendMessage(msg, "Your Pomodoro timer is started")
					go a.users[userName].Run()

					continue
				}

				if _, ok := a.users[userName]; ok {
					a.SendMessage(msg, "Your Pomodoro timer is stopped")
					a.users[userName].Stop()
					delete(a.users, userName)
					continue
				}

				a.SendMessage(msg, "You don't have a Pomodoro timer launched")
			}
		}
	}
}

func (a App) SendMessage(msg *tgbotapi.Message, msgText string) {
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Your Pomodoro timer is started")
	newMsg.ReplyToMessageID = msg.MessageID
	if _, err := a.bot.Send(newMsg); err != nil {
		log.Println("Could not send message:", err)
	}
}
