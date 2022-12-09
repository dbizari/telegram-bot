package game_stages

import user_pkg "tdl/internal/domain/user"

type Police struct {
}

func (p Police) GetStageName() string {
	return STAGE_POLICE
}

func (p Police) CanUserVote(user user_pkg.UserInfo) bool {
	return false
}
