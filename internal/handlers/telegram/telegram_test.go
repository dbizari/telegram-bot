package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"tdl/internal/handlers/cmd"
	mock_cmd "tdl/testing/mocks/handlers_mock/cmd"
	mock_getter "tdl/testing/mocks/handlers_mock/cmdgetter"
	mock_telegram "tdl/testing/mocks/handlers_mock/telegram"
	"testing"
)

func TestTelegramHandler_HandleUpdate(t *testing.T) {
	tests := []struct {
		name            string
		incomingMsg     tgbotapi.Update
		wantErr         bool
		msgErr          string
		fnMockCmdGetter func(cmdGetterMock *mock_getter.MockCmdGetter, cmdMock *mock_cmd.MockCmdHandler)
		fnMockBotAPI    func(botMock *mock_telegram.MockBotAPI)
	}{
		{
			name: "Happy path",
			incomingMsg: tgbotapi.Update{
				Message: &tgbotapi.Message{
					MessageID: 55667,
					Chat:      &tgbotapi.Chat{ID: 123},
					From:      &tgbotapi.User{UserName: "el_dani_pa"},
				},
			},
			wantErr: false,
			fnMockCmdGetter: func(cmdGetterMock *mock_getter.MockCmdGetter, cmdMock *mock_cmd.MockCmdHandler) {
				expectedPayload := cmd.CmdPayload{
					Args:     []string{"arg1", "arg2"},
					UserName: "el_dani_pa",
					ChatID:   123,
				}
				cmdMock.EXPECT().HandleCmd(gomock.Any(), expectedPayload).Times(1).
					Return("command success", nil)
				cmdGetterMock.EXPECT().GetCmdAndArgsFromMessage(gomock.Any()).Times(1).
					Return(cmdMock, []string{"arg1", "arg2"})
			},
			fnMockBotAPI: func(botMock *mock_telegram.MockBotAPI) {
				botMock.EXPECT().SendMsg(int64(123), "command success", 55667).Times(1).Return(nil)
			},
		},
		{
			name: "invalid cmd",
			incomingMsg: tgbotapi.Update{
				Message: &tgbotapi.Message{
					MessageID: 55667,
					Chat:      &tgbotapi.Chat{ID: 123},
					From:      &tgbotapi.User{UserName: "el_dani_pa"},
				},
			},
			wantErr: true,
			msgErr:  "invalid command",
			fnMockCmdGetter: func(cmdGetterMock *mock_getter.MockCmdGetter, cmdMock *mock_cmd.MockCmdHandler) {
				cmdGetterMock.EXPECT().GetCmdAndArgsFromMessage(gomock.Any()).Times(1).
					Return(nil, nil)
			},
			fnMockBotAPI: func(botMock *mock_telegram.MockBotAPI) {
				botMock.EXPECT().SendMsg(int64(123), "invalid command", 55667).Times(1).Return(nil)
			},
		},
		{
			name: "HandleCmd returns with error",
			incomingMsg: tgbotapi.Update{
				Message: &tgbotapi.Message{
					MessageID: 55667,
					Chat:      &tgbotapi.Chat{ID: 123},
					From:      &tgbotapi.User{UserName: "el_dani_pa"},
				},
			},
			wantErr: true,
			msgErr:  "something went wrong command failed :(",
			fnMockCmdGetter: func(cmdGetterMock *mock_getter.MockCmdGetter, cmdMock *mock_cmd.MockCmdHandler) {
				expectedPayload := cmd.CmdPayload{
					Args:     []string{"arg1", "arg2"},
					UserName: "el_dani_pa",
					ChatID:   123,
				}
				cmdMock.EXPECT().HandleCmd(gomock.Any(), expectedPayload).Times(1).
					Return("", errors.New("command failed :("))
				cmdGetterMock.EXPECT().GetCmdAndArgsFromMessage(gomock.Any()).Times(1).
					Return(cmdMock, []string{"arg1", "arg2"})
			},
			fnMockBotAPI: func(botMock *mock_telegram.MockBotAPI) {
				botMock.EXPECT().SendMsg(int64(123), "something went wrong command failed :(", 55667).Times(1).Return(nil)
			},
		},
		{
			name: "Success operation but Send message fails",
			incomingMsg: tgbotapi.Update{
				Message: &tgbotapi.Message{
					MessageID: 55667,
					Chat:      &tgbotapi.Chat{ID: 123},
					From:      &tgbotapi.User{UserName: "el_dani_pa"},
				},
			},
			wantErr: false,
			fnMockCmdGetter: func(cmdGetterMock *mock_getter.MockCmdGetter, cmdMock *mock_cmd.MockCmdHandler) {
				expectedPayload := cmd.CmdPayload{
					Args:     []string{"arg1", "arg2"},
					UserName: "el_dani_pa",
					ChatID:   123,
				}
				cmdMock.EXPECT().HandleCmd(gomock.Any(), expectedPayload).Times(1).
					Return("command success", nil)
				cmdGetterMock.EXPECT().GetCmdAndArgsFromMessage(gomock.Any()).Times(1).
					Return(cmdMock, []string{"arg1", "arg2"})
			},
			fnMockBotAPI: func(botMock *mock_telegram.MockBotAPI) {
				botMock.EXPECT().SendMsg(int64(123), "command success", 55667).Times(1).Return(errors.New("telegram API failed"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			botMock := mock_telegram.NewMockBotAPI(ctrl)
			cmdGetterMock := mock_getter.NewMockCmdGetter(ctrl)
			cmdMock := mock_cmd.NewMockCmdHandler(ctrl)

			tt.fnMockBotAPI(botMock)
			tt.fnMockCmdGetter(cmdGetterMock, cmdMock)

			th := TelegramHandler{
				BotAPI:    botMock,
				CmdGetter: cmdGetterMock,
			}

			err := th.HandleUpdate(tt.incomingMsg)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.msgErr != "" {
				assert.Contains(t, tt.msgErr, err.Error())
			}
		})
	}
}
