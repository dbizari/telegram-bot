package create_game

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"tdl/internal/domain/game_session"
	user_pkg "tdl/internal/domain/user"
	"tdl/internal/handlers/cmd"
	mock_repository "tdl/testing/mocks/repository_mock"
	"testing"
)

func TestCreateGameSessionHandler_HandleCmd(t *testing.T) {
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
					Args:     []string{"randomArg"},
					UserName: "danybiz",
				},
			},
			want:    fmt.Sprintf(REPLY_CREATE_GAME, "create-game-id"),
			wantErr: false,
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				repository.EXPECT().CreateGame(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, gameSession *game_session.GameSession) (string, error) {
						if !gameSession.ID.IsZero() {
							return "", fmt.Errorf("expected gameSession.ID to be 0")
						}

						expectedStatus := "pending"
						if gameSession.Stage.GetStageName() != expectedStatus {
							return "", fmt.Errorf("expected status: %s , received: %s", expectedStatus, gameSession.Stage)
						}

						expectedUsername := "danybiz"
						if gameSession.OwnerId != expectedUsername {
							return "", fmt.Errorf("expected username: %s , received: %s", expectedUsername, gameSession.OwnerId)
						}

						expectedUsers := []*user_pkg.UserInfo{{
							UserId:   "danybiz",
							Role:     "",
							Alive:    true,
							Votes:    0,
							HasVoted: false,
						}}
						if len(gameSession.Users) != 1 {
							return "", fmt.Errorf("expected len(gameSession.Users): %v, received: %v", len(expectedUsers), len(gameSession.Users))
						}

						if !reflect.DeepEqual(gameSession.Users[0], expectedUsers[0]) {
							return "", fmt.Errorf("expected userInfo: %+v, received: %+v", expectedUsers[0], gameSession.Users[0])
						}

						return "create-game-id", nil
					})
			},
		},
		{
			name: "userName is missing on gameSession",
			args: args{
				payload: cmd.CmdPayload{
					Args:     []string{"randomArg"},
					UserName: "",
				},
			},
			wantErr: true,
			msgErr:  "error on create game session handler, username should not be empty",
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				repository.EXPECT().CreateGame(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "repository fails to create the game",
			args: args{
				payload: cmd.CmdPayload{
					Args:     []string{"randomArg"},
					UserName: "danybiz",
				},
			},
			wantErr: true,
			msgErr:  "unexpected error with mongo",
			fnMockRepository: func(repository *mock_repository.MockGameSessionRepositoryAPI) {
				repository.EXPECT().CreateGame(gomock.Any(), gomock.Any()).Times(1).
					Return("", errors.New("unexpected error with mongo"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repositoryMock := mock_repository.NewMockGameSessionRepositoryAPI(ctrl)
			tt.fnMockRepository(repositoryMock)

			handler := CreateGameSessionHandler{
				Repository: repositoryMock,
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
