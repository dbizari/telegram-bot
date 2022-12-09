package create_game

import (
	"context"
	"fmt"
	"tdl/internal/domain/game_session"
	"tdl/internal/domain/game_stages"
	user_pkg "tdl/internal/domain/user"
	"tdl/internal/handlers/cmd"
	"tdl/internal/repository"
)

const (
	CMD_CREATE_GAME   = "/createGame"
	REPLY_CREATE_GAME = "Game created with id: %s. Share it with your friends !"
)

type CreateGameSessionHandler struct {
	Repository repository.GameSessionRepositoryAPI
}

func (cgsh CreateGameSessionHandler) HandleCmd(ctx context.Context, payload cmd.CmdPayload) (string, error) {
	if payload.UserName == "" {
		return "", fmt.Errorf("error on create game session handler, username should not be empty")
	}

	gameSession := &game_session.GameSession{
		OwnerId: payload.UserName,
		Users: []*user_pkg.UserInfo{
			{
				UserId:   payload.UserName,
				Role:     "", // ToDo for the moment this is empty, previous start the game the role should be assigned
				Alive:    true,
				Votes:    0,
				HasVoted: false,
			},
		},
		Stage: game_stages.Pending{},
	}

	id, err := cgsh.Repository.CreateGame(ctx, gameSession)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(REPLY_CREATE_GAME, id), nil
}
