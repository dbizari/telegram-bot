package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	// Game Stages
	STAGE_PENDING    = "pending"
	STAGE_MAFIA      = "mafia"
	STAGE_POLICE     = "police"
	STAGE_DISCUSSION = "discussion"
	STAGE_FINISHIED  = "finished"

	// Game Roles
	ROLE_MAFIA   = "mafia"
	ROLE_CITIZEN = "citizen"
	ROLE_POLICE  = "police"
)

type GameSession struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	OwnerId string             `json:"owner_id" bson:"owner_id"`
	Users   []*UserInfo        `json:"users" bson:"users"`
	Status  string             `json:"status" bson:"status"`
}

type UserInfo struct {
	UserId   string `json:"user_id" bson:"user_id"`
	Role     string `json:"role" bson:"role"`
	Alive    bool   `json:"alive" bson:"alive"`
	Votes    int    `json:"votes" bson:"votes"`
	HasVoted bool   `json:"has_voted" bson:"has_voted"`
}

func (gs GameSession) CanUserVote(userID string) bool {
	var user *UserInfo
	for _, u := range gs.Users {
		if u.UserId == userID {
			user = u
		}
	}

	if !user.Alive {
		return false
	}

	if gs.Status == STAGE_DISCUSSION {
		return !user.HasVoted
	}

	if gs.Status == STAGE_MAFIA {
		return user.Role == ROLE_MAFIA && !user.HasVoted
	}

	return false
}

func (gs *GameSession) ApplyVote(votingUserID, votedUserID string) bool {
	found := false
	for _, user := range gs.Users {
		if user.UserId == votingUserID {
			user.HasVoted = true
		}
		if user.UserId == votedUserID {
			user.Votes++
			found = true
		}
	}

	return found
}

func (gs GameSession) GetRole(userId string) string {
	var role string

	for _, user := range gs.Users {
		if user.UserId == userId {
			role = user.Role
		}
	}

	return role
}

func (gs GameSession) CanUserAskForRole(userId string, userToAsk string) bool {
	userRole := gs.GetRole(userId)

	if userRole == "" {
		return false
	}

	if userId == userToAsk || (userRole == ROLE_POLICE && gs.Status == STAGE_POLICE) {
		return true
	}

	return false
}
