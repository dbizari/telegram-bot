package game_stages

import (
	"fmt"
	user_pkg "tdl/internal/domain/user"
)

func BuildVotationList(users []string, action string) string {
	msg := fmt.Sprintf("You have to select one user to %s:\n", action)

	for i, user := range users {
		msg += fmt.Sprintf("%d. %s\n", i, user)
	}

	return msg
}

func getMostVotedUser(users []*user_pkg.UserInfo) *user_pkg.UserInfo {
	max := 0
	var votedUser *user_pkg.UserInfo
	for _, u := range users {
		if u.Votes >= max {
			max = u.Votes
			votedUser = u
		}
	}

	return votedUser
}
