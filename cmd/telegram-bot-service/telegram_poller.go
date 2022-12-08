package main

import (
	"log"
	"tdl/internal/client"
	"tdl/internal/handlers/cmd/getter"
	"tdl/internal/handlers/telegram"
)

func startTelegramPoller() {
	bot := client.GetTelegramBotClient()

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
