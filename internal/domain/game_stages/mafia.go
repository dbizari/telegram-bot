package game_stages

import (
	"tdl/internal/clients/telegram"
	user_pkg "tdl/internal/domain/user"
)

type Mafia struct {
}

func (m Mafia) IsVotationDone(users []*user_pkg.UserInfo) bool {
	for _, u := range users {
		if u.Role == user_pkg.ROLE_MAFIA && u.HasVoted == false {
			return false
		}
	}

	return true
}

func (m Mafia) ApplyAction(users []*user_pkg.UserInfo) {
	//TODO implement, recuento de votos, funar al mas votado (alive = false) y hacer broadcast de quien murio
	panic("implement me")
}

func (m Mafia) NextStage(users []*user_pkg.UserInfo) GameStage {
	// ToDo mandar mensaje al policia de que elija...
	return Police{}
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
