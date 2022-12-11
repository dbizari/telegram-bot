package join_game

import (
	"context"
	"fmt"
	user_pkg "tdl/internal/domain/user"
	"tdl/internal/handlers/cmd"
	"tdl/internal/repository"
)

const (
	CMD_JOIN_GAME                    = "/joinGame"
	REPLY_JOIN_GAME                  = "Joined to game: %s!"
	REPLY_JOIN_ALREADY_ON_OTHER_GAME = "You are already on other game"
	REPLY_JOIN_INEXISTENT_SESSION    = "The game you are trying to join does not exist"
)

type JoinGameSessionHandler struct {
	Repository repository.GameSessionRepositoryAPI
}

func (jgsh JoinGameSessionHandler) HandleCmd(ctx context.Context, payload cmd.CmdPayload) (string, error) {
	if payload.UserName == "" {
		return "", fmt.Errorf("error on join game session handler, username should not be empty")
	}
	if len(payload.Args) == 0 {
		return "", fmt.Errorf("error: missing session ID")
	}
	sessionId := payload.Args[0]

	gameSession, err := jgsh.Repository.GetNotFinishedGameByMember(ctx, payload.UserName)
	if err != nil {
		return "", err
	}
	if gameSession != nil && !gameSession.ID.IsZero() {
		return REPLY_JOIN_ALREADY_ON_OTHER_GAME, nil
	}

	gameSession, err = jgsh.Repository.Get(ctx, sessionId)
	if err != nil {
		return "", err
	}
	if gameSession == nil {
		return REPLY_JOIN_INEXISTENT_SESSION, nil
	}

	newUser := &user_pkg.UserInfo{
		UserId:   payload.UserName,
		ChatID:   payload.ChatID,
		Role:     "",
		Alive:    true,
		Votes:    0,
		HasVoted: false,
	}

	gameSession.Users = append(gameSession.Users, newUser)

	err = jgsh.Repository.Update(ctx, gameSession)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(REPLY_JOIN_GAME, gameSession.ID.Hex()), nil
}
