package handlers

import (
	"context"
	"strings"
	"tdl/internal/repository"
)

type CmdHandler interface {
	HandleCmd(ctx context.Context, payload CmdPayload) (string, error)
}

type CmdPayload struct {
	Args     []string
	UserName string
}

func GetCmdAndArgsFromMessage(message string) (CmdHandler, []string) {
	splittedMessage := strings.Split(message, " ")

	// Match command
	var cmd CmdHandler
	switch splittedMessage[0] {
	case CMD_CREATE_GAME:
		cmd = &CreateGameSessionHandler{Repository: repository.GetGameSessionRepositoryClient()}
	case CMD_JOIN_GAME:
		cmd = &JoinGameSessionHandler{Repository: repository.GetGameSessionRepositoryClient()}
	default:
		// unrecognizable command
		return nil, nil
	}

	return cmd, splittedMessage[1:]
}
