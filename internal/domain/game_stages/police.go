package game_stages

import (
	"tdl/internal/clients/telegram"
	user_pkg "tdl/internal/domain/user"
)

type Police struct {
}

func (p Police) IsVotationDone(users []*user_pkg.UserInfo) bool {
	return false
}

func (p Police) ApplyAction(users []*user_pkg.UserInfo) {
	// Esto se resuelve con el comando /ask role
	panic("shouldn't be called")
}

func (p Police) NextStage(users []*user_pkg.UserInfo) GameStage {
	chatIDs := make([]int64, 0)
	aliveUsers := make([]string, 0)
	for _, u := range users {
		if u.Alive {
			aliveUsers = append(aliveUsers, u.UserId)
			chatIDs = append(chatIDs, u.ChatID)
		}
	}

	telegram.GetTelegramBotClient().BroadcastMsgToUsers(chatIDs, "Discussion time... Feel free to chat with the players")
	telegram.GetTelegramBotClient().BroadcastMsgToUsers(chatIDs, BuildVotationList(aliveUsers, "kick"))

	return Discussion{}
}

func (p Police) GetStageName() string {
	return STAGE_POLICE
}

func (p Police) CanUserVote(user user_pkg.UserInfo) bool {
	return false
}

func (p Police) Start(users []*user_pkg.UserInfo) {
	policeChatIDs := make([]int64, 0)
	nonPoliceUsers := make([]string, 0)
	for _, user := range users {
		if user.Role == user_pkg.ROLE_POLICE {
			policeChatIDs = append(policeChatIDs, user.ChatID)
		} else {
			nonPoliceUsers = append(nonPoliceUsers, user.UserId)
		}
	}

	telegram.GetTelegramBotClient().
		BroadcastMsgToUsers(policeChatIDs, BuildVotationList(nonPoliceUsers, "ask role"))
}
