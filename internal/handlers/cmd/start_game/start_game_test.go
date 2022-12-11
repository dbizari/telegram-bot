package start_game

import (
	"context"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tdl/internal/clients/telegram"
	"tdl/internal/domain/game_session"
	"tdl/internal/domain/game_stages"
	user_pkg "tdl/internal/domain/user"
	"tdl/internal/handlers/cmd"
	mock_telegram "tdl/testing/mocks/handlers_mock/telegram"
	mock_repository "tdl/testing/mocks/repository_mock"
	"testing"
)

func TestStartGameHandler_HandleCmd(t *testing.T) {
	tests := []struct {
		name              string
		fnMockRepository  func(repository *mock_repository.MockGameSessionRepositoryAPI)
		fnMockTelegramBot func(mock *mock_telegram.MockBotAPI)
		args              cmd.CmdPayload
		want              string
		wantErr           bool
	}{
		{
			name: "Happy path",
			args: cmd.CmdPayload{
				UserName: "mili",
			},
			want:    "",
			wantErr: false,
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := game_session.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "mili",
					Users: []*user_pkg.UserInfo{
						{
							UserId: "mili",
							ChatID: 1,
						},
						{
							UserId: "danybiz",
							ChatID: 2,
						},
						{
							UserId: "tfanciotti",
							ChatID: 3,
						},
					},
					Stage: game_stages.Pending{},
				}
				repository.EXPECT().GetNotFinishedGameByMember(gomock.Any(), "mili").Times(1).
					Return(&session, nil)

				repository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).
					Return(nil)
			},
			fnMockTelegramBot: func(mock *mock_telegram.MockBotAPI) {
				mock.EXPECT().BroadcastMsgToUsers(gomock.Any(), gomock.Any()).Times(2)
				mock.EXPECT().SendMsg(gomock.Any(), gomock.Any(), gomock.Any()).Times(3)
			},
		},
		{
			name: "User is not the owner",
			args: cmd.CmdPayload{
				UserName: "danybiz",
			},
			want:    REPLY_START_GAME_USER_IS_NOT_OWNER,
			wantErr: false,
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := game_session.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "mili",
					Users: []*user_pkg.UserInfo{
						{
							UserId: "mili",
						},
						{
							UserId: "danybiz",
						},
						{
							UserId: "tfanciotti",
						},
					},
					Stage: game_stages.Mafia{},
				}
				repository.EXPECT().GetNotFinishedGameByMember(gomock.Any(), "danybiz").Times(1).
					Return(&session, nil)
			},
			fnMockTelegramBot: func(mock *mock_telegram.MockBotAPI) {},
		},
		{
			name: "Game already started",
			args: cmd.CmdPayload{
				UserName: "mili",
			},
			want:    REPLY_START_GAME_ALREADY_STARTED,
			wantErr: false,
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := game_session.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "mili",
					Users: []*user_pkg.UserInfo{
						{
							UserId: "mili",
						},
						{
							UserId: "danybiz",
						},
						{
							UserId: "tfanciotti",
						},
					},
					Stage: game_stages.Mafia{},
				}
				repository.EXPECT().GetNotFinishedGameByMember(gomock.Any(), "mili").Times(1).
					Return(&session, nil)
			},
			fnMockTelegramBot: func(mock *mock_telegram.MockBotAPI) {},
		},
		{
			name: "Not enough players",
			args: cmd.CmdPayload{
				UserName: "mili",
			},
			want:    REPLY_START_GAME_NOT_ENOUGH_PLAYERS,
			wantErr: false,
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := game_session.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "mili",
					Users: []*user_pkg.UserInfo{
						{
							UserId: "mili",
						},
						{
							UserId: "danybiz",
						},
					},
					Stage: game_stages.Pending{},
				}
				repository.EXPECT().GetNotFinishedGameByMember(gomock.Any(), "mili").Times(1).
					Return(&session, nil)
			},
			fnMockTelegramBot: func(mock *mock_telegram.MockBotAPI) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repositoryMock := mock_repository.NewMockGameSessionRepositoryAPI(ctrl)
			tt.fnMockRepository(repositoryMock)

			tbMock := mock_telegram.NewMockBotAPI(ctrl)
			tt.fnMockTelegramBot(tbMock)
			telegram.SetMockTelegramBot(tbMock)

			handler := StartGameHandler{
				GameSessionRepository: repositoryMock,
			}

			got, err := handler.HandleCmd(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("HandleCmd() got = %v, want %v", got, tt.want)
			}
		})
	}
}
