package ask_role

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

func TestAskRoleHandler_HandleCmd(t *testing.T) {
	tests := []struct {
		name              string
		fnMockRepository  func(repository *mock_repository.MockGameSessionRepositoryAPI)
		fnMockTelegramBot func(mock *mock_telegram.MockBotAPI)
		args              cmd.CmdPayload
		want              string
		wantErr           bool
	}{
		{
			name: "Happy path with police",
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := game_session.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "danybiz",
					Users: []*user_pkg.UserInfo{
						{
							UserId: "danybiz",
							Role:   user_pkg.ROLE_POLICE,
						},
						{
							UserId: "tomi",
							Role:   user_pkg.ROLE_MAFIA,
						},
					},
					Stage: game_stages.Police{},
				}
				repository.EXPECT().GetNotFinishedGameByMember(gomock.Any(), "danybiz").Times(1).
					Return(&session, nil)
				repository.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1)
			},
			fnMockTelegramBot: func(mock *mock_telegram.MockBotAPI) {
				mock.EXPECT().BroadcastMsgToUsers(gomock.Any(), gomock.Any()).Times(2)
			},
			args: cmd.CmdPayload{
				UserName: "danybiz",
				Args:     []string{"tomi"},
			},
			want:    user_pkg.ROLE_MAFIA,
			wantErr: false,
		},
		{
			name: "Missing username",
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
			},
			fnMockTelegramBot: func(mock *mock_telegram.MockBotAPI) {},
			args: cmd.CmdPayload{
				UserName: "danybiz",
			},
			want:    REPLY_ASK_ROLE_MISSING_USERNAME,
			wantErr: false,
		},
		{
			name: "Citizen can't ask for another user's role",
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := game_session.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "danybiz",
					Users: []*user_pkg.UserInfo{
						{
							UserId: "danybiz",
							Role:   user_pkg.ROLE_CITIZEN,
						},
						{
							UserId: "tomi",
							Role:   user_pkg.ROLE_MAFIA,
						},
					},
					Stage: game_stages.Discussion{},
				}
				repository.EXPECT().GetNotFinishedGameByMember(gomock.Any(), "danybiz").Times(1).
					Return(&session, nil)
			},
			fnMockTelegramBot: func(mock *mock_telegram.MockBotAPI) {},
			args: cmd.CmdPayload{
				UserName: "danybiz",
				Args:     []string{"tomi"},
			},
			want:    REPLY_ASK_ROLE_USER_CANT_KNOW_ROLE,
			wantErr: false,
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

			handler := AskRoleHandler{
				GameSessionRepository: repositoryMock,
			}
			got, err := handler.HandleCmd(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HandleCmd() got = %v, want %v", got, tt.want)
			}
		})
	}
}
