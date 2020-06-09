package main

import (
	"errors"
	"os"

	"github.com/niakr1s/chatty-server/app/bot"
	log "github.com/sirupsen/logrus"
)

func main() {
	e, err := parseEnv()
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := b.Connect(e.BotUsername, e.BotPassword, e.URL); err != nil {
		log.Fatal(err)
	}
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
