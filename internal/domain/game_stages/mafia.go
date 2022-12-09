package game_stages

import (
	"tdl/internal/clients/telegram"
	user_pkg "tdl/internal/domain/user"
)

type Mafia struct {
}

func (m Mafia) CanUserVote(user user_pkg.UserInfo) bool {
	return user.Role == user_pkg.ROLE_MAFIA && !user.HasVoted
}

func (m Mafia) GetStageName() string {
	return STAGE_MAFIA
}

func (m Mafia) Start(users []*user_pkg.UserInfo) {
	mafiaChatIDs := make([]int64, 0)
	nonMafiaUsers := make([]string, 0)
	for _, user := range users {
		if user.Role == user_pkg.ROLE_MAFIA {
			mafiaChatIDs = append(mafiaChatIDs, user.ChatID)
		} else {
			nonMafiaUsers = append(nonMafiaUsers, user.UserId)
		}
	}

	telegram.GetTelegramBotClient().
		BroadcastMsgToUsers(mafiaChatIDs, BuildVotationList(nonMafiaUsers, "kill"))
}
