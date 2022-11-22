package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"tdl/internal/handlers/cmd/getter"
	"tdl/internal/handlers/telegram"
)

func startTelegramPoller() {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
		return
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	telegramHandler := telegram.TelegramHandler{
		BotAPI:    bot,
		CmdGetter: getter.CmdGetterImpl{},
	}

	log.Println("Poller successfuly initiated...")

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			go telegramHandler.HandleUpdate(update)
		}
	}
}
