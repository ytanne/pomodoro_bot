package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCommand   = "/start"
	finishCommand  = "/finish"
	setRestTime    = "/set_rest_time"
	setWorkingTime = "/set_work_time"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDoesNotExist  = errors.New("user does not exist")
)

type App struct {
	bot   *tgbotapi.BotAPI
	users map[int64]worker
	mutex *sync.Mutex
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
		tgbotapi.BotCommand{
			Command:     setRestTime,
			Description: "Set rest time (in minutes)",
		},
		tgbotapi.BotCommand{
			Command:     setWorkingTime,
			Description: "Set working time (in minutes)",
		},
	)

	if _, err := bot.Request(setCommands); err != nil {
		return App{}, fmt.Errorf("could not set up commands: %w", err)
	}

	mu := sync.Mutex{}

	return App{
		bot:   bot,
		mutex: &mu,
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

				switch msg.Command() {
				case startCommand[1:]:
					newWorker := worker{
						bot:            a.bot,
						chatID:         msg.Chat.ID,
						ctx:            make(chan struct{}),
						chillingPeriod: defaultChillingPeriod,
						hustlingPeriod: defaultHustlingPeriod,
					}

					if err := a.AddNewUser(userID, newWorker); err != nil {
						a.SendMessage(msg, "Your Pomodoro timer is already working")
						continue
					}

					if err := a.LaunchTimer(userID); err != nil {
						a.SendMessage(msg, "Could not start Pomodoro timer. Contact administrator about potential bug")
						return
					}

					a.SendMessage(msg, "Your Pomodoro timer has started! Now go to work")
				case finishCommand[1:]:
					if err := a.StopTimer(userID); err != nil {
						a.SendMessage(msg, "You don't have a Pomodoro timer launched")
						continue
					}

					a.SendMessage(msg, "Your Pomodoro timer has stopped. See you next time!")
				case setRestTime[1:]:
					chillingPeriod, err := retrieveInputDuration(msg.CommandArguments())
					if errors.Is(err, ErrInvalidInput) {
						a.SendMessage(msg, "The input is invalid, please enter a valid number")
						continue
					} else if err != nil {
						a.SendMessage(msg, err.Error())
						continue
					}

					if err := a.SetupChillingDuration(userID, time.Duration(chillingPeriod)); err != nil {
						a.SendMessage(msg, "You don't have a Pomodoro timer launched")
						continue
					}
					a.SendMessage(msg, fmt.Sprintf("Your rest time is updated to %d minutes", chillingPeriod))
				case setWorkingTime[1:]:
					hustlingTimePeriod, err := retrieveInputDuration(msg.CommandArguments())
					if errors.Is(err, ErrInvalidInput) {
						a.SendMessage(msg, "The input is invalid, please enter a valid number")
						continue
					} else if err != nil {
						a.SendMessage(msg, err.Error())
						continue
					}

					if err := a.SetupHustlingDuration(userID, time.Duration(hustlingTimePeriod)); err != nil {
						a.SendMessage(msg, "You don't have a Pomodoro timer launched")
						continue
					}
					a.SendMessage(msg, fmt.Sprintf("Your work time is updated to %d minutes", hustlingTimePeriod))
				default:
					a.SendMessage(msg, "unknow command obtained - %s", msg.Command())
				}
			}
		}
	}
}

func retrieveInputDuration(arg string) (int, error) {
	if arg == "" {
		return -1, errors.New("input must be a valid integer from 1 to 90")
	}

	value, err := strconv.Atoi(arg)
	if err != nil {
		return -1, fmt.Errorf("entered value - %s; error: %w", arg, ErrInvalidInput)
	}

	if value <= 0 || value > 90 {
		return -1, fmt.Errorf("entered value - %s; error: %w", arg, ErrInvalidInput)
	}

	return value, nil
}

func (a App) SendMessage(msg *tgbotapi.Message, format string, values ...interface{}) {
	msgText := fmt.Sprintf(format, values...)
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	newMsg.ReplyToMessageID = msg.MessageID
	if _, err := a.bot.Send(newMsg); err != nil {
		log.Println("Could not send message:", err)
	}
}
