package getter

import (
	"strings"
	"tdl/internal/handlers/cmd"
	"tdl/internal/handlers/cmd/create_game"
	"tdl/internal/handlers/cmd/start_game"
	"tdl/internal/handlers/cmd/vote"
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
	case start_game.CMD_START_GAME:
		command = &start_game.StartGameHandler{GameSessionRepository: repository.GetGameSessionRepositoryClient()}
	case vote.CMD_VOTE:
		command = &vote.VoteHandler{GameSessionRepository: repository.GetGameSessionRepositoryClient()}
	default:
		// unrecognizable command
		return nil, nil
	}

	return command, splittedMessage[1:]
}
