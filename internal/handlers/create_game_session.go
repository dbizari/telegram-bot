package handlers

import (
	"context"
)

const CMD_CREATE_GAME = "/createGame"

type CreateGameSessionHandler struct {
}

func (gsh CreateGameSessionHandler) HandleCmd(ctx context.Context, payload CmdPayload) (string, error) {
	return "", nil
}
