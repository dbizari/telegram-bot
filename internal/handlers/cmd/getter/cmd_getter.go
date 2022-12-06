package getter

import (
	"strings"
	"tdl/internal/handlers/cmd"
	"tdl/internal/handlers/cmd/create_game"
	"tdl/internal/handlers/cmd/exit_game"
	"tdl/internal/handlers/cmd/join_game"
	"tdl/internal/repository"
)

type CmdGetter interface {
	GetCmdAndArgsFromMessage(message string) (cmd.CmdHandler, []string)
}

type CmdGetterImpl struct {
}

func (cgi CmdGetterImpl) GetCmdAndArgsFromMessage(message string) (cmd.CmdHandler, []string) {
	splittedMessage := strings.Split(message, " ")

	// Match command
	var command cmd.CmdHandler
	switch splittedMessage[0] {
	case create_game.CMD_CREATE_GAME:
		command = &create_game.CreateGameSessionHandler{Repository: repository.GetGameSessionRepositoryClient()}
	case join_game.CMD_JOIN_GAME:
		command = &join_game.JoinGameSessionHandler{Repository: repository.GetGameSessionRepositoryClient()}
	case exit_game.CMD_EXIT_GAME:
		command = &exit_game.ExitGameSessionHandler{Repository: repository.GetGameSessionRepositoryClient()}
	default:
		// unrecognizable command
		return nil, nil
	}

	return command, splittedMessage[1:]
}
