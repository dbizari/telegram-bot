package game_stages

import (
	"fmt"
	"tdl/internal/clients/telegram"
	user_pkg "tdl/internal/domain/user"
)

type Discussion struct {
}

func (d Discussion) IsVotingDone(users []*user_pkg.UserInfo) bool {
	for _, u := range users {
		if u.Alive && u.HasVoted == false {
			return false
		}
	}

	return true
}

func (d Discussion) ApplyAction(users []*user_pkg.UserInfo) {
	votedUser := getMostVotedUser(users)
	votedUser.Alive = false
	telegram.GetTelegramBotClient().SendMsg(votedUser.ChatID, "you were voted to be kicked out of the game :(", 0)

	chatIDs := make([]int64, 0)
	for _, u := range users {
		if u.Alive {
			chatIDs = append(chatIDs, u.ChatID)
		}
	}

	telegram.GetTelegramBotClient().BroadcastMsgToUsers(chatIDs, fmt.Sprintf("%s was kicked out with role %s", votedUser.UserId, votedUser.Role))
}

func (d Discussion) NextStage(users []*user_pkg.UserInfo) GameStage {
	mafiaCount := 0
	citizenCount := 0

	chatIDs := make([]int64, 0)
	for _, u := range users {
		chatIDs = append(chatIDs, u.ChatID)

		if !u.Alive {
			continue
		}

		if u.Role == user_pkg.ROLE_MAFIA {
			mafiaCount++
		} else {
			citizenCount++
		}
	}

	if mafiaCount == 0 {
		telegram.GetTelegramBotClient().BroadcastMsgToUsers(chatIDs, "Citizens wins !")
		return Finished{}
	}

	if mafiaCount > citizenCount || (len(users) == 2 && mafiaCount == 1) {
		telegram.GetTelegramBotClient().BroadcastMsgToUsers(chatIDs, "Mafia wins !")
		return Finished{}
	}

	nextStage := Mafia{}
	nextStage.Start(users)
	return Mafia{}
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
		if !user.Alive {
			continue
		}
		chatIDs = append(chatIDs, user.ChatID)
		allUsers = append(allUsers, user.UserId)
	}

	telegram.GetTelegramBotClient().
		BroadcastMsgToUsers(chatIDs, "Debating time! Who should be kicked out of the city?")
	telegram.GetTelegramBotClient().
		BroadcastMsgToUsers(chatIDs, BuildVotationList(allUsers, "kick"))
}
