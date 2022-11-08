package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TelegramHandler struct {
	BotAPI *tgbotapi.BotAPI
}

func (th TelegramHandler) HandleUpdate(incomingMsg tgbotapi.Update) {
	cmd, args := GetCmdAndArgsFromMessage(incomingMsg.Message.Text)
	if cmd == nil {
		th.sendReplyMsg(incomingMsg, "invalid command")
		return
	}

	reply, err := cmd.HandleCmd(context.Background(), CmdPayload{
		Args:     args,
		UserName: incomingMsg.Message.From.UserName,
	})
	if err != nil {
		// handle error
		th.sendReplyMsg(incomingMsg, "something went wrong"+err.Error())
		return
	}

	th.sendReplyMsg(incomingMsg, reply)
}

func (th TelegramHandler) sendReplyMsg(incomingMsg tgbotapi.Update, reply string) {
	msg := tgbotapi.NewMessage(incomingMsg.Message.Chat.ID, reply)
	msg.ReplyToMessageID = incomingMsg.Message.MessageID

	if _, err := th.BotAPI.Send(msg); err != nil {
		log.Printf("error replying message to %s", incomingMsg.Message.From.UserName)
	}
}
