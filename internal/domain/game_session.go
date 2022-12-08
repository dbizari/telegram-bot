package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

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

	// Game restrictions regarding players
	RESTRICTION_GAME_MIN_PLAYERS                       = 3
	RESTRICTION_PLAYERS_AMOUNT_NEED_MORE_SPECIAL_ROLES = 6
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

func (gs GameSession) CanUserStartTheGame() bool {
	return gs.Status == STAGE_PENDING
}

func (gs *GameSession) StartGame() bool {
	usersAmount := len(gs.Users)

	if usersAmount < RESTRICTION_GAME_MIN_PLAYERS {
		return false
	}

	// The game starts with 3 players: 1 mafia, 1 police and 1 citizen
	// Every time 6 new players are added, one police and one mafia will be added
	specialRolesAmount := 1 + (usersAmount-RESTRICTION_GAME_MIN_PLAYERS)/RESTRICTION_PLAYERS_AMOUNT_NEED_MORE_SPECIAL_ROLES
	policeUsers := 0
	mafiaUsers := 0

	// Initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	for policeUsers < specialRolesAmount {
		randomPos := rand.Intn(usersAmount)
		user := gs.Users[randomPos]
		if user.Role == ROLE_POLICE {
			continue
		}
		user.Role = ROLE_POLICE
		policeUsers++
	}

	for mafiaUsers < specialRolesAmount {
		randomPos := rand.Intn(usersAmount)
		user := gs.Users[randomPos]
		if user.Role == ROLE_POLICE || user.Role == ROLE_MAFIA {
			continue
		}
		user.Role = ROLE_MAFIA
		mafiaUsers++
	}

	for _, user := range gs.Users {
		if user.Role == ROLE_MAFIA || user.Role == ROLE_POLICE {
			continue
		}
		user.Role = ROLE_CITIZEN
	}

	gs.Status = STAGE_MAFIA

	return true
}

func (gs GameSession) IsUserTheOwner(userId string) bool {
	return gs.OwnerId == userId
}
