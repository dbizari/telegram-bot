package handlers

import (
	"context"
	"fmt"
	"tdl/internal/domain"
	"tdl/internal/repository"
)

const CMD_CREATE_GAME = "/createGame"

type CreateGameSessionHandler struct {
	Repository repository.GameSessionRepositoryAPI
}

func (cgsh CreateGameSessionHandler) HandleCmd(ctx context.Context, payload CmdPayload) (string, error) {
	gameSession := &domain.GameSession{
		OwnerId: payload.UserName,
		Users: []domain.UserInfo{
			{
				UserId: payload.UserName,
				Role:   "", // ToDo for the moment this is empty, previous start the game the role should be assigned
				Alive:  true,
			},
		},
		Status: "pending",
	}

	id, err := cgsh.Repository.CreateGame(ctx, gameSession)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Game created with id: %s. Share it with your friends !", id), nil
}
