package telegram

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tdl/internal/client"
	"tdl/internal/handlers/cmd"
	"tdl/internal/handlers/cmd/getter"
)

type TelegramHandler struct {
	client.BotAPI
	getter.CmdGetter
}

func (th TelegramHandler) HandleUpdate(incomingMsg tgbotapi.Update) error {
	command, args := th.GetCmdAndArgsFromMessage(incomingMsg.Message.Text)
	if command == nil {
		if err := th.SendMsg(incomingMsg.Message.Chat.ID, "invalid command", incomingMsg.Message.MessageID); err != nil {
			log.Printf("error replying message to %s", incomingMsg.Message.From.UserName)
		}
		return errors.New("invalid command")
	}

	reply, err := command.HandleCmd(context.Background(), cmd.CmdPayload{
		Args:     args,
		UserName: incomingMsg.Message.From.UserName,
		ChatID:   incomingMsg.Message.Chat.ID,
	})
	if err != nil {
		// handle error
		msgErr := "something went wrong " + err.Error()
		if err := th.SendMsg(incomingMsg.Message.Chat.ID, msgErr, incomingMsg.Message.MessageID); err != nil {
			log.Printf("error replying message to %s", incomingMsg.Message.From.UserName)
		}
		return err
	}

	if err := th.SendMsg(incomingMsg.Message.Chat.ID, reply, incomingMsg.Message.MessageID); err != nil {
		log.Printf("error replying message to %s", incomingMsg.Message.From.UserName)
	}

	return nil
}
