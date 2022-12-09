package game_stages

import (
	user_pkg "tdl/internal/domain/user"
)

type Pending struct {
}

func (p Pending) CanUserVote(user user_pkg.UserInfo) bool {
	return false
}

func (p Pending) GetStageName() string {
	return STAGE_PENDING
}

func (p Pending) Start(users []*user_pkg.UserInfo) {
	// Nothing to de here
}
