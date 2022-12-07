package vote

import (
	"context"
	"fmt"
	"tdl/internal/handlers/cmd"
	"tdl/internal/repository"
)

const (
	CMD_VOTE                      = "/vote"
	REPLY_VOTE                    = "voted username: %s"
	REPLY_VOTE_MISSING_USERNAME   = "missing username to vote"
	REPLY_VOTE_INEXISTENT_SESSION = "oops it seems to be that you are not in a game"
	REPLY_VOTE_USER_CANT_VOTE     = "sorry you cannot vote :("
	REPLY_VOTE_USER_NOT_FOUND     = "%s not found in the game"
)

type VoteHandler struct {
	GameSessionRepository repository.GameSessionRepositoryAPI
}

func (vh *VoteHandler) HandleCmd(ctx context.Context, payload cmd.CmdPayload) (string, error) {
	if len(payload.Args) != 1 {
		return REPLY_VOTE_MISSING_USERNAME, nil
	}
	if payload.Args[0] == "" {
		return REPLY_VOTE_MISSING_USERNAME, nil
	}

	session, err := vh.GameSessionRepository.GetByMember(ctx, payload.UserName)
	if err != nil {
		return "", err
	}
	if session == nil {
		return REPLY_VOTE_INEXISTENT_SESSION, nil
	}

	if !session.CanUserVote(payload.UserName) {
		return REPLY_VOTE_USER_CANT_VOTE, nil
	}

	ok := session.ApplyVote(payload.UserName, payload.Args[0])
	if !ok {
		return fmt.Sprintf(REPLY_VOTE_USER_NOT_FOUND, payload.Args[0]), nil
	}

	err = vh.GameSessionRepository.Update(ctx, session)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(REPLY_VOTE, payload.Args[0]), nil
}
