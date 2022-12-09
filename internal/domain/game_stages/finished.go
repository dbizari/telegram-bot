package game_stages

import user_pkg "tdl/internal/domain/user"

type Finished struct {
}

func (f Finished) IsVotationDone(users []*user_pkg.UserInfo) bool {
	//TODO implement me
	panic("implement me")
}

func (f Finished) ApplyAction(users []*user_pkg.UserInfo) {
	//TODO implement me
	panic("implement me")
}

func (f Finished) NextStage(users []*user_pkg.UserInfo) GameStage {
	//TODO implement me
	panic("implement me")
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
