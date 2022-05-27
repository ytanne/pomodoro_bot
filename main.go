package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ytanne/pomodoro_bot/pkg/app"
)

func main() {
	token := flag.String("token", "", "Token of Telegram bot")
	flag.Parse()

	program, err := app.NewApp(*token)
	if err != nil {
		log.Println("Could not create new Pomodoro program:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	log.Println("Launching the program...")

	program.Run(ctx)

	log.Println("Program finished...")
}
