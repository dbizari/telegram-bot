package game_session

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	_ "tdl/internal/clients/telegram"
	"tdl/internal/domain/game_stages"
	user_pkg "tdl/internal/domain/user"
	"time"
)

const (
	// Game restrictions regarding players
	RESTRICTION_GAME_MIN_PLAYERS                       = 3
	RESTRICTION_PLAYERS_AMOUNT_NEED_MORE_SPECIAL_ROLES = 6
)

type GameSession struct {
	ID      primitive.ObjectID    `json:"_id" bson:"_id,omitempty"`
	OwnerId string                `json:"owner_id" bson:"owner_id"`
	Users   []*user_pkg.UserInfo  `json:"users" bson:"users"`
	Stage   game_stages.GameStage `json:"status" bson:"-"` // this field will be custom-marshaled
}

func (gs GameSession) MarshalBSON() ([]byte, error) {
	type GameSessionMirrorBSON struct {
		GameSession GameSession `bson:",inline"`
		Status      interface{} `bson:"status"`
	}
	return bson.Marshal(GameSessionMirrorBSON{
		GameSession: gs,
		Status:      gs.Stage.GetStageName(),
	})
}

func (gs *GameSession) UnmarshalBSON(data []byte) error {
	type GameSessionMirrorBSON struct {
		GameSession GameSession `bson:",inline"`
		Status      interface{} `bson:"status"`
	}
	var gameSessionMirror GameSessionMirrorBSON
	err := bson.Unmarshal(data, &gameSessionMirror)
	if err != nil {
		return err
	}

	gs.ID = gameSessionMirror.GameSession.ID
	gs.Users = gameSessionMirror.GameSession.Users
	gs.OwnerId = gameSessionMirror.GameSession.OwnerId

	switch gameSessionMirror.Status {
	case game_stages.STAGE_PENDING:
		gs.Stage = game_stages.Pending{}
	case game_stages.STAGE_MAFIA:
		gs.Stage = game_stages.Mafia{}
	case game_stages.STAGE_POLICE:
		gs.Stage = game_stages.Police{}
	case game_stages.STAGE_DISCUSSION:
		gs.Stage = game_stages.Discussion{}
	case game_stages.STAGE_FINISHIED:
		gs.Stage = game_stages.Finished{}
	}

	return nil
}

func (gs GameSession) CanUserVote(userID string) bool {
	var user *user_pkg.UserInfo
	for _, u := range gs.Users {
		if u.UserId == userID {
			user = u
		}
	}

	if !user.Alive {
		return false
	}

	return gs.Stage.CanUserVote(*user)
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
	return gs.Stage.GetStageName() == game_stages.STAGE_PENDING
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
		if user.Role == user_pkg.ROLE_POLICE {
			continue
		}
		user.Role = user_pkg.ROLE_POLICE
		policeUsers++
	}

	for mafiaUsers < specialRolesAmount {
		randomPos := rand.Intn(usersAmount)
		user := gs.Users[randomPos]
		if user.Role == user_pkg.ROLE_POLICE || user.Role == user_pkg.ROLE_MAFIA {
			continue
		}
		user.Role = user_pkg.ROLE_MAFIA
		mafiaUsers++
	}

	for _, user := range gs.Users {
		if user.Role == user_pkg.ROLE_MAFIA || user.Role == user_pkg.ROLE_POLICE {
			continue
		}
		user.Role = user_pkg.ROLE_CITIZEN
	}

	gs.Stage = game_stages.Mafia{}

	return true
}

func (gs GameSession) IsUserTheOwner(userId string) bool {
	return gs.OwnerId == userId
}

func (gs GameSession) StartStage() {
	gs.Stage.Start(gs.Users)
}

func (gs *GameSession) ApplyStageAction() {
	if gs.Stage.IsVotationDone(gs.Users) {
		gs.Stage.ApplyAction(gs.Users)
		gs.Stage = gs.Stage.NextStage(gs.Users)

		for _, u := range gs.Users {
			u.HasVoted = false
			u.Votes = 0
		}
	}
}
