package start_game

import (
	"context"
	"tdl/internal/handlers/cmd"
	"tdl/internal/repository"
)

const (
	CMD_START_GAME = "/startGame"

	REPLY_START_GAME                    = "Let the game begin!"
	REPLY_START_GAME_INEXISTENT_SESSION = "Oops it seems that you are not in any game"
	REPLY_START_GAME_USER_NOT_THE_OWNER = "Sorry, only the owner can start the game"
	REPLY_START_GAME_NOT_ENOUGH_PLAYERS = "More players are needed to start the game!"
)

type StartGameHandler struct {
	GameSessionRepository repository.GameSessionRepositoryAPI
}

func (sgh StartGameHandler) HandleCmd(ctx context.Context, payload cmd.CmdPayload) (string, error) {
	session, err := sgh.GameSessionRepository.GetByMember(ctx, payload.UserName)
	if err != nil {
		return "", err
	}

	if session == nil {
		return REPLY_START_GAME_INEXISTENT_SESSION, nil
	}

	if !session.CanUserStartTheGame(payload.UserName) {
		return REPLY_START_GAME_USER_NOT_THE_OWNER, nil
	}

	ok := session.StartGame()
	if !ok {
		return REPLY_START_GAME_NOT_ENOUGH_PLAYERS, nil
	}

	err = sgh.GameSessionRepository.Update(ctx, session)
	if err != nil {
		return "", err
	}

	return REPLY_START_GAME, nil
}
