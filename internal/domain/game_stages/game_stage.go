package game_stages

import user_pkg "tdl/internal/domain/user"

const (
	// Game Stages, the flow is pending -> mafia -> police -> discussion -> finished
	STAGE_PENDING    = "pending"
	STAGE_MAFIA      = "mafia"
	STAGE_POLICE     = "police"
	STAGE_DISCUSSION = "discussion"
	STAGE_FINISHIED  = "finished"
)

type GameStage interface {
	GetStageName() string
	CanUserVote(user user_pkg.UserInfo) bool
	Start(users []*user_pkg.UserInfo)
	IsVotationDone(users []*user_pkg.UserInfo) bool
	ApplyAction(users []*user_pkg.UserInfo)
	NextStage(users []*user_pkg.UserInfo) GameStage
}
