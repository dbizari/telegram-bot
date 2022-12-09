package game_stages

import (
	"fmt"
	user_pkg "tdl/internal/domain/user"
)

func BuildVotationList(users []string, action string) string {
	msg := fmt.Sprintf("You have to select 1 user to %s:\n", action)

	for i, user := range users {
		msg += fmt.Sprintf("%d. %s\n", i, user)
	}

	return msg
}

func getMostVotedUser(users []*user_pkg.UserInfo) *user_pkg.UserInfo {
	votes := make(map[*user_pkg.UserInfo]int)

	for _, u := range users {
		votes[u]++
	}

	max := 0
	var votedUser *user_pkg.UserInfo
	for k, v := range votes {
		if v >= max {
			max = v
			votedUser = k
		}
	}

	return votedUser
}
