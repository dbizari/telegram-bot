package user

const (
	// Game Roles
	ROLE_MAFIA   = "mafia"
	ROLE_CITIZEN = "citizen"
	ROLE_POLICE  = "police"
)

type UserInfo struct {
	UserId   string `json:"user_id" bson:"user_id"`
	ChatID   int64  `json:"chat_id" bson:"chat_id"`
	Role     string `json:"role" bson:"role"`
	Alive    bool   `json:"alive" bson:"alive"`
	Votes    int    `json:"votes" bson:"votes"`
	HasVoted bool   `json:"has_voted" bson:"has_voted"`
}
