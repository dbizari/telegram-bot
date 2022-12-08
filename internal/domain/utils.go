package domain

import "fmt"

func buildVotationList(users []string, action string) string {
	msg := fmt.Sprintf("You have to select 1 user to %s:\n", action)

	for i, user := range users {
		msg += fmt.Sprintf("%d. %s\n", i, user)
	}

	return msg
}
