package ask_role

import (
	"context"
	"tdl/internal/handlers/cmd"
	"tdl/internal/repository"
)

const (
	CMD_ASK_ROLE                        = "/askRole"
	REPLY_ASK_ROLE_MISSING_USERNAME     = "Missing username to know their role"
	REPLY_ASK_ROLE_INEXISTENT_SESSION   = "Oops it seems that you are not in a game"
	REPLY_ASK_ROLE_USER_CANT_KNOW_ROLE  = "You can't ask for that user's role"
	REPLY_ASK_ROLE_USER_NOT_IN_THE_GAME = "This user is not playing the game"
)

type AskRoleHandler struct {
	GameSessionRepository repository.GameSessionRepositoryAPI
}

func (arh *AskRoleHandler) HandleCmd(ctx context.Context, payload cmd.CmdPayload) (string, error) {
	if len(payload.Args) != 1 {
		return REPLY_ASK_ROLE_MISSING_USERNAME, nil
	}
	if payload.Args[0] == "" {
		return REPLY_ASK_ROLE_MISSING_USERNAME, nil
	}

	session, err := arh.GameSessionRepository.GetNotFinishedGameByMember(ctx, payload.UserName)
	if err != nil {
		return "", err
	}
	if session == nil {
		return REPLY_ASK_ROLE_INEXISTENT_SESSION, nil
	}

	if !session.CanUserAskForRole(payload.UserName) {
		return REPLY_ASK_ROLE_USER_CANT_KNOW_ROLE, nil
	}

	role := session.GetRole(payload.Args[0])
	if role == "" {
		return REPLY_ASK_ROLE_USER_NOT_IN_THE_GAME, nil
	}

	session.Stage = session.Stage.NextStage(session.Users)

	err = arh.GameSessionRepository.Update(ctx, session)
	if err != nil {
		return "", err
	}

	return role, nil
}
