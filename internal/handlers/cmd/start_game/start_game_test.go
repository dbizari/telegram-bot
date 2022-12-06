package start_game

import (
	"context"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tdl/internal/domain"
	"tdl/internal/handlers/cmd"
	mock_repository "tdl/testing/mocks/repository_mock"
	"testing"
)

func TestStartGameHandler_HandleCmd(t *testing.T) {
	tests := []struct {
		name             string
		fnMockRepository func(repository *mock_repository.MockGameSessionRepositoryAPI)
		args             cmd.CmdPayload
		want             string
		wantErr          bool
	}{
		{
			name: "Happy path",
			args: cmd.CmdPayload{
				UserName: "mili",
			},
			want:    REPLY_START_GAME,
			wantErr: false,
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := domain.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "mili",
					Users: []*domain.UserInfo{
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
					Status: domain.STAGE_PENDING,
				}
				repository.EXPECT().GetByMember(gomock.Any(), "mili").Times(1).
					Return(&session, nil)

				expectedSession := session
				expectedSession.Status = domain.STAGE_MAFIA
				expectedSession.Users[0].Role = domain.ROLE_POLICE
				expectedSession.Users[1].Role = domain.ROLE_MAFIA
				expectedSession.Users[2].Role = domain.ROLE_CITIZEN

				repository.EXPECT().Update(gomock.Any(), gomock.Eq(&expectedSession)).Times(1).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repositoryMock := mock_repository.NewMockGameSessionRepositoryAPI(ctrl)
			tt.fnMockRepository(repositoryMock)

			handler := StartGameHandler{
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
