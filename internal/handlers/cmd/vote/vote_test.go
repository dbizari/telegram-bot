package vote

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tdl/internal/domain"
	"tdl/internal/handlers/cmd"
	mock_repository "tdl/testing/mocks/repository_mock"
	"testing"
)

func TestVoteHandler_HandleCmd(t *testing.T) {
	type args struct {
		payload cmd.CmdPayload
	}
	tests := []struct {
		name             string
		args             args
		want             string
		wantErr          bool
		msgErr           string
		fnMockRepository func(repository *mock_repository.MockGameSessionRepositoryAPI)
	}{
		{
			name: "Happy path",
			args: args{
				payload: cmd.CmdPayload{
					Args:     []string{"tfanciotti"},
					UserName: "danybiz",
				},
			},
			want:    fmt.Sprintf(REPLY_VOTE, "tfanciotti"),
			wantErr: false,
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				session := domain.GameSession{
					ID:      primitive.ObjectID{},
					OwnerId: "danybiz",
					Users: []*domain.UserInfo{
						{
							UserId:   "danybiz",
							Role:     domain.ROLE_MAFIA,
							Alive:    true,
							Votes:    0,
							HasVoted: false,
						},
						{
							UserId:   "tfanciotti",
							Role:     domain.ROLE_CITIZEN,
							Alive:    true,
							Votes:    0,
							HasVoted: false,
						},
					},
					Status: domain.STAGE_DISCUSSION,
				}
				repository.EXPECT().GetByMember(gomock.Any(), "danybiz").Times(1).
					Return(&session, nil)

				var expectedSession domain.GameSession = session
				expectedSession.Users[0].HasVoted = false
				expectedSession.Users[1].Votes++

				repository.EXPECT().Update(gomock.Any(), gomock.Eq(&expectedSession)).Times(1).
					Return(nil)
			},
		},
		// ToDo terminar de cubrir los casos
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repositoryMock := mock_repository.NewMockGameSessionRepositoryAPI(ctrl)
			tt.fnMockRepository(repositoryMock)

			handler := VoteHandler{
				GameSessionRepository: repositoryMock,
			}
			got, err := handler.HandleCmd(context.Background(), tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.msgErr != "" {
				assert.Contains(t, tt.msgErr, err.Error())
			}
			if got != tt.want {
				t.Errorf("HandleCmd() got = %v, want %v", got, tt.want)
			}
		})
	}
}
