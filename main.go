package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ytanne/pomodoro_bot/pkg/app"
	"github.com/ytanne/pomodoro_bot/pkg/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Config file to run the bot")
	flag.Parse()

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Println("Could not get config:", err)
		os.Exit(1)
	}

	program, err := app.NewApp(cfg)
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
