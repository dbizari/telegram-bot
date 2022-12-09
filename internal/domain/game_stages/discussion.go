package game_stages

import user_pkg "tdl/internal/domain/user"

type Discussion struct {
}

func (d Discussion) CanUserVote(user user_pkg.UserInfo) bool {
	return !user.HasVoted
}

func (d Discussion) GetStageName() string {
	return STAGE_DISCUSSION
}
