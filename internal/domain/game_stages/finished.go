package game_stages

import user_pkg "tdl/internal/domain/user"

type Finished struct {
}

func (f Finished) CanUserVote(user user_pkg.UserInfo) bool {
	return false
}

func (f Finished) GetStageName() string {
	return STAGE_FINISHIED
}

func (f Finished) Start(users []*user_pkg.UserInfo) {
	panic("IMPLEMENT ME")
}
