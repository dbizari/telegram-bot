package game_stages

import (
	user_pkg "tdl/internal/domain/user"
)

type Pending struct {
}

func (p Pending) IsVotationDone(users []*user_pkg.UserInfo) bool {
	return false
}

func (p Pending) ApplyAction(users []*user_pkg.UserInfo) {
	// ToDo estaria copado mover lo de mili aca
	panic("implement me")
}

func (p Pending) NextStage(users []*user_pkg.UserInfo) GameStage {
	return Mafia{}
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
