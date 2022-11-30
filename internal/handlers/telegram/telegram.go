package telegram

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tdl/internal/handlers/cmd"
	"tdl/internal/handlers/cmd/getter"
)

type TelegramHandler struct {
	BotAPI
	getter.CmdGetter
}

// Interface for the tgbotapi.BotAPI in order to be able to do dependency injection for mock
type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func (th TelegramHandler) HandleUpdate(incomingMsg tgbotapi.Update) error {
	command, args := th.GetCmdAndArgsFromMessage(incomingMsg.Message.Text)
	if command == nil {
		th.sendReplyMsg(incomingMsg, "invalid command")
		return errors.New("invalid command")
	}

	reply, err := command.HandleCmd(context.Background(), cmd.CmdPayload{
		Args:     args,
		UserName: incomingMsg.Message.From.UserName,
	})
	if err != nil {
		// handle error
		th.sendReplyMsg(incomingMsg, "something went wrong "+err.Error())
		return err
	}

	th.sendReplyMsg(incomingMsg, reply)
	return nil
}

func (th TelegramHandler) sendReplyMsg(incomingMsg tgbotapi.Update, reply string) {
	msg := tgbotapi.NewMessage(incomingMsg.Message.Chat.ID, reply)
	msg.ReplyToMessageID = incomingMsg.Message.MessageID

	if _, err := th.BotAPI.Send(msg); err != nil {
		log.Printf("error replying message to %s", incomingMsg.Message.From.UserName)
	}
}
