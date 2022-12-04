package exit_game

import (
	"context"
	"fmt"
	"tdl/internal/handlers/cmd"
	"tdl/internal/repository"
)

const (
	CMD_EXIT_GAME   = "/exitGame"
	REPLY_EXIT_GAME = "Bye!"
)

type ExitGameSessionHandler struct {
	Repository repository.GameSessionRepositoryAPI
}

func (egsh ExitGameSessionHandler) HandleCmd(ctx context.Context, payload cmd.CmdPayload) (string, error) {
	if payload.UserName == "" {
		return "", fmt.Errorf("error on exit game session handler, username should not be empty")
	}

	result, err := egsh.Repository.ExitGame(ctx, payload.UserName)
	if err != nil {
		return "", err
	}

	if result == false {
		return "", fmt.Errorf("you are not in any game")
	}

	return fmt.Sprintf(REPLY_EXIT_GAME), nil
}
