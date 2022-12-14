package game_stages

import (
	"fmt"
	"tdl/internal/clients/telegram"
	user_pkg "tdl/internal/domain/user"
)

type Mafia struct {
}

func (m Mafia) IsVotingDone(users []*user_pkg.UserInfo) bool {
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

	telegram.GetTelegramBotClient().SendMsg(votedUser.ChatID, "You were killed by the mafia...", 0)

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
	isAnyPoliceAlive := false
	var usersAlive int

	for _, u := range users {
		chatIDs = append(chatIDs, u.ChatID)
		if u.Role == user_pkg.ROLE_POLICE && u.Alive {
			isAnyPoliceAlive = true
		}
		if u.Alive {
			usersAlive++
		}
	}

	// 1 mafia and other or 2 mafias
	if usersAlive == 2 {
		telegram.GetTelegramBotClient().BroadcastMsgToUsers(chatIDs, "Mafia wins !")
		return Finished{}
	}

	var nextStage GameStage

	if isAnyPoliceAlive {
		nextStage = Police{}
	} else {
		nextStage = Discussion{}
	}

	nextStage.Start(users)

	return nextStage
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
	nonMafiaChatIDs := make([]int64, 0)
	for _, user := range users {
		if !user.Alive {
			continue
		}

		if user.Role == user_pkg.ROLE_MAFIA {
			mafiaChatIDs = append(mafiaChatIDs, user.ChatID)
		} else {
			nonMafiaUsers = append(nonMafiaUsers, user.UserId)
			nonMafiaChatIDs = append(nonMafiaChatIDs, user.ChatID)
		}
	}

	telegram.GetTelegramBotClient().
		BroadcastMsgToUsers(nonMafiaChatIDs, "The mafia is up to something...")

	telegram.GetTelegramBotClient().
		BroadcastMsgToUsers(mafiaChatIDs, BuildVotationList(nonMafiaUsers, "kill"))
}
