package ask_role

import (
	"context"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tdl/internal/domain"
	"tdl/internal/handlers/cmd"
	mock_repository "tdl/testing/mocks/repository_mock"
	"testing"
)

func TestAskRoleHandler_HandleCmd(t *testing.T) {
	tests := []struct {
		name             string
		fnMockRepository func(repository *mock_repository.MockGameSessionRepositoryAPI)
		args             cmd.CmdPayload
		want             string
		wantErr          bool
	}{
		{
			name: "Happy path with police",
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := domain.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "danybiz",
					Users: []*domain.UserInfo{
						{
							UserId: "danybiz",
							Role:   domain.ROLE_POLICE,
						},
						{
							UserId: "tomi",
							Role:   domain.ROLE_MAFIA,
						},
					},
					Status: domain.STAGE_POLICE,
				}
				repository.EXPECT().GetByMember(gomock.Any(), "danybiz").Times(1).
					Return(&session, nil)
			},
			args: cmd.CmdPayload{
				UserName: "danybiz",
				Args:     []string{"tomi"},
			},
			want:    domain.ROLE_MAFIA,
			wantErr: false,
		},
		{
			name: "Happy path with user",
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := domain.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "danybiz",
					Users: []*domain.UserInfo{
						{
							UserId: "danybiz",
							Role:   domain.ROLE_CITIZEN,
						},
						{
							UserId: "tomi",
							Role:   domain.ROLE_MAFIA,
						},
					},
					Status: domain.STAGE_DISCUSSION,
				}
				repository.EXPECT().GetByMember(gomock.Any(), "danybiz").Times(1).
					Return(&session, nil)
			},
			args: cmd.CmdPayload{
				UserName: "danybiz",
				Args:     []string{"danybiz"},
			},
			want:    domain.ROLE_CITIZEN,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repositoryMock := mock_repository.NewMockGameSessionRepositoryAPI(ctrl)
			tt.fnMockRepository(repositoryMock)

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
