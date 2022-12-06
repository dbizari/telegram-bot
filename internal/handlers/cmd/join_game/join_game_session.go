package join_game

import (
	"context"
	"fmt"
	"tdl/internal/domain"
	"tdl/internal/handlers/cmd"
	"tdl/internal/repository"
)

const (
	CMD_JOIN_GAME   = "/joinGame"
	REPLY_JOIN_GAME = "Joined to game: %s!"
)

type JoinGameSessionHandler struct {
	Repository repository.GameSessionRepositoryAPI
}

func (jgsh JoinGameSessionHandler) HandleCmd(ctx context.Context, payload cmd.CmdPayload) (string, error) {
	if payload.UserName == "" {
		return "", fmt.Errorf("error on join game session handler, username should not be empty")
	}

	newUser := &domain.UserInfo{
		UserId: payload.UserName,
		Role:   "", // ToDo for the moment this is empty, previous start the game the role should be assigned
		Alive:  true,
	}

	if len(payload.Args) == 0 {
		return "", fmt.Errorf("error: missing session ID")
	}

	sesionId := payload.Args[0]

	id, err := jgsh.Repository.AddPlayer(ctx, sesionId, newUser)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(REPLY_JOIN_GAME, id), nil
}
