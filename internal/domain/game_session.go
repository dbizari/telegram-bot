package domain

type GameSession struct {
	Id      string     `json:"id"`
	OwnerId string     `json:"owner_id"`
	Users   []UserInfo `json:"users"`
	Status  string     `json:"status"`
}

type UserInfo struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	Alive  bool   `json:"alive"`
}
