package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"log"
	"os"
	"sync"
)

var (
	once     sync.Once
	instance BotAPI
)

// Interface for the tgbotapi.BotAPI in order to be able to do dependency injection for mock
type BotAPI interface {
	GetUpdatesChan() tgbotapi.UpdatesChannel
	SendMsg(chatID int64, msg string, messageToReplyID int) error
	BroadcastMsgToUsers(chatIDs []int64, msg string)
}

type telegramBotImpl struct {
	Bot *tgbotapi.BotAPI
}

func GetTelegramBotClient() BotAPI {
	once.Do(func() {
		telegramToken := os.Getenv("TELEGRAM_TOKEN")
		bot, err := tgbotapi.NewBotAPI(telegramToken)
		if err != nil {
			log.Panic(err)
			return
		}
		bot.Debug = true
		log.Printf("Authorized on account %s", bot.Self.UserName)

		instance = telegramBotImpl{
			Bot: bot,
		}
	})

	return instance
}

func SetMockTelegramBot(mock BotAPI) {
	once.Do(func() {})
	instance = mock
}

func (tb telegramBotImpl) GetUpdatesChan() tgbotapi.UpdatesChannel {
	// Setup updates channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return tb.Bot.GetUpdatesChan(u)
}

func (tb telegramBotImpl) BroadcastMsgToUsers(chatIDs []int64, msg string) {
	if len(chatIDs) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	for _, chat := range chatIDs {
		wg.Add(1)
		go func(c int64) {
			err := tb.SendMsg(c, msg, 0)
			if err != nil {
				log.Printf("error sending message to %d", c)
			}
			wg.Done()
		}(chat)
	}

	wg.Wait()
}

func (tb telegramBotImpl) SendMsg(chatID int64, msg string, messageToReplyID int) error {
	payload := tgbotapi.NewMessage(chatID, msg)
	if messageToReplyID != 0 {
		payload.ReplyToMessageID = messageToReplyID
	}

	if _, err := tb.Bot.Send(payload); err != nil {
		return errors.Wrap(err, "error sending message")
	}

	return nil
}
