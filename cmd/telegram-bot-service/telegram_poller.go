package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"tdl/internal/handlers"
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

	log.Println("Poller successfuly initiated...")

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			cmd, args := handlers.GetCmdAndArgsFromMessage(update.Message.Text)
			if cmd == nil {
				sendReplyMsg(bot, update, "invalid command")
			}

			reply, err := cmd.HandleCmd(context.Background(), handlers.CmdPayload{
				Args:     args,
				UserName: update.Message.From.UserName,
			})
			if err != nil {
				sendReplyMsg(bot, update, "something went wrong"+err.Error())
			}

			sendReplyMsg(bot, update, reply)
		}
	}
}

func sendReplyMsg(bot *tgbotapi.BotAPI, incomingMsg tgbotapi.Update, reply string) {
	msg := tgbotapi.NewMessage(incomingMsg.Message.Chat.ID, reply)
	msg.ReplyToMessageID = incomingMsg.Message.MessageID

	if _, err := bot.Send(msg); err != nil {
		log.Printf("error replying message to %s", incomingMsg.Message.From.UserName)
	}
}
