package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/niakr1s/chatty-server/app/bot"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.TraceLevel)

	e, err := parseEnv()
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.NewHelloBot()
	if err != nil {
		log.Fatal(err)
	}

	if err := b.Connect(e.BotUsername, e.BotPassword, e.URL); err != nil {
		log.Fatal(err)
	}

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)
	signal.Notify(exit, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		b.Run(ctx)
		done <- struct{}{}
	}()

	<-exit
	cancel()
	<-done
}

type env struct {
	BotUsername string
	BotPassword string
	URL         string
}

func parseEnv() (env, error) {
	res := env{}

	res.BotUsername = os.Getenv("BOT_USERNAME")
	if res.BotUsername == "" {
		return res, errors.New("empty BOT_USERNAME")
	}

	res.BotPassword = os.Getenv("BOT_PASSWORD")
	if res.BotPassword == "" {
		return res, errors.New("empty BOT_PASSWORD")
	}

	res.URL = os.Getenv("BOT_URL")
	if res.URL == "" {
		return res, errors.New("empty BOT_URL")
	}

	return res, nil
}