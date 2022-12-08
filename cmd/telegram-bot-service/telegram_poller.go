package main

import (
	"log"
	telegram_client "tdl/internal/clients/telegram"
	"tdl/internal/handlers/cmd/getter"
	"tdl/internal/handlers/telegram"
)

func startTelegramPoller() {
	bot := telegram_client.GetTelegramBotClient()

	telegramHandler := telegram.TelegramHandler{
		BotAPI:    bot,
		CmdGetter: getter.CmdGetterImpl{},
	}

	log.Println("Poller successfuly initiated...")

	for update := range bot.GetUpdatesChan() {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			go telegramHandler.HandleUpdate(update)
		}
	}
}
