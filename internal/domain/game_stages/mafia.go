package game_stages

import user_pkg "tdl/internal/domain/user"

type Mafia struct {
}

func (m Mafia) CanUserVote(user user_pkg.UserInfo) bool {
	return user.Role == user_pkg.ROLE_MAFIA && !user.HasVoted
}

func (m Mafia) GetStageName() string {
	return STAGE_MAFIA
}
