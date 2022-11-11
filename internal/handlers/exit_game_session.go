package handlers

import (
	"context"
	"fmt"
	"tdl/internal/repository"
)

const (
	CMD_EXIT_GAME   = "/exitGame"
	REPLY_EXIT_GAME = "Bye!"
)

type ExitGameSessionHandler struct {
	Repository repository.GameSessionRepositoryAPI
}

func (cgsh ExitGameSessionHandler) HandleCmd(ctx context.Context, payload CmdPayload) (string, error) {
	if payload.UserName == "" {
		return "", fmt.Errorf("error on create game session handler, username should not be empty")
	}

	result, err := cgsh.Repository.ExitGame(ctx, payload.UserName)
	if err != nil {
		return "", err
	}

	if result == false {
		return "", fmt.Errorf("You are not in any game.")
	}

	return fmt.Sprintf(REPLY_EXIT_GAME), nil
}
