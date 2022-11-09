package handlers

import (
	"context"
	"fmt"
	"tdl/internal/domain"
	"tdl/internal/repository"
)

const (
	CMD_JOIN_GAME   = "/joinGame"
	REPLY_JOIN_GAME = "Joined to game: %s!"
)

type JoinGameSessionHandler struct {
	Repository repository.GameSessionRepositoryAPI
}

func (cgsh JoinGameSessionHandler) HandleCmd(ctx context.Context, payload CmdPayload) (string, error) {
	if payload.UserName == "" {
		return "", fmt.Errorf("error on join game session handler, username should not be empty")
	}

	newUser := &domain.UserInfo{
		UserId: payload.UserName,
		Role:   "", // ToDo for the moment this is empty, previous start the game the role should be assigned
		Alive:  true,
	}

	if len(payload.Args) == 0 {
		return "", fmt.Errorf("Error: missing session ID")
	}

	sesionId := payload.Args[0]

	id, err := cgsh.Repository.AddPlayer(ctx, sesionId, newUser)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(REPLY_JOIN_GAME, id), nil
}
