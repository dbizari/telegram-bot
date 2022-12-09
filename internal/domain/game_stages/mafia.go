package game_stages

import (
	"fmt"
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
	votedUser := getMostVotedUser(users)
	votedUser.Alive = false

	telegram.GetTelegramBotClient().SendMsg(votedUser.ChatID, "you was killed by the mafia...", 0)

	chatIDs := make([]int64, 0)
	for _, u := range users {
		if u.Alive {
			chatIDs = append(chatIDs, u.ChatID)
		}
	}

	telegram.GetTelegramBotClient().BroadcastMsgToUsers(chatIDs, fmt.Sprintf("Unfortunately %s was killed in this round...", votedUser.UserId))
}

func (m Mafia) NextStage(users []*user_pkg.UserInfo) GameStage {
	chatIDs := make([]int64, 0)
	aliveUsers := make([]string, 0)
	for _, u := range users {
		if u.Role == user_pkg.ROLE_POLICE {
			chatIDs = append(chatIDs, u.ChatID)
		} else {
			if u.Alive {
				aliveUsers = append(aliveUsers, u.UserId)
			}
		}
	}

	telegram.GetTelegramBotClient().BroadcastMsgToUsers(chatIDs, BuildVotationList(aliveUsers, "ask role"))

	// ToDo si no hay policia vivo se tiene que pasar directo a discussion

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
		if !user.Alive {
			continue
		}

		if user.Role == user_pkg.ROLE_MAFIA {
			mafiaChatIDs = append(mafiaChatIDs, user.ChatID)
		} else {
			nonMafiaUsers = append(nonMafiaUsers, user.UserId)
		}
	}

	telegram.GetTelegramBotClient().
		BroadcastMsgToUsers(mafiaChatIDs, BuildVotationList(nonMafiaUsers, "kill"))
}
