package game_stages

import (
	"tdl/internal/clients/telegram"
	user_pkg "tdl/internal/domain/user"
)

type Discussion struct {
}

func (d Discussion) CanUserVote(user user_pkg.UserInfo) bool {
	return !user.HasVoted
}

func (d Discussion) GetStageName() string {
	return STAGE_DISCUSSION
}

func (d Discussion) Start(users []*user_pkg.UserInfo) {
	chatIDs := make([]int64, 0)
	allUsers := make([]string, 0)
	for _, user := range users {
		chatIDs = append(chatIDs, user.ChatID)
		allUsers = append(allUsers, user.UserId)
	}

	telegram.GetTelegramBotClient().
		BroadcastMsgToUsers(chatIDs, BuildVotationList(allUsers, "kick"))
}
