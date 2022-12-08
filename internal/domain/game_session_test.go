package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

func TestGameSession_CanUserVote(t *testing.T) {
	type fields struct {
		ID      primitive.ObjectID
		OwnerId string
		Users   []*UserInfo
		Status  string
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Mafia person votes on Mafia stage",
			fields: fields{
				Users: []*UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     ROLE_MAFIA,
					HasVoted: false,
				}},
				Status: STAGE_MAFIA,
			},
			args: args{
				userID: "dani",
			},
			want: true,
		},
		{
			name: "Dead user cannot vote",
			fields: fields{
				Users: []*UserInfo{{
					UserId: "dani",
					Alive:  false,
				}},
			},
			args: args{
				userID: "dani",
			},
			want: false,
		},
		{
			name: "Citizen cannot vote on Mafia stage",
			fields: fields{
				Users: []*UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     ROLE_CITIZEN,
					HasVoted: false,
				}},
				Status: STAGE_MAFIA,
			},
			args: args{
				userID: "dani",
			},
			want: false,
		},
		{
			name: "Mafia person votes on Discussion stage",
			fields: fields{
				Users: []*UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     ROLE_MAFIA,
					HasVoted: false,
				}},
				Status: STAGE_DISCUSSION,
			},
			args: args{
				userID: "dani",
			},
			want: true,
		},
		{
			name: "Citizen person votes on Discussion stage",
			fields: fields{
				Users: []*UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     ROLE_MAFIA,
					HasVoted: false,
				}},
				Status: STAGE_DISCUSSION,
			},
			args: args{
				userID: "dani",
			},
			want: true,
		},
		{
			name: "Citizen person votes on Police stage",
			fields: fields{
				Users: []*UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     ROLE_MAFIA,
					HasVoted: false,
				}},
				Status: STAGE_POLICE,
			},
			args: args{
				userID: "dani",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameSession{
				ID:      tt.fields.ID,
				OwnerId: tt.fields.OwnerId,
				Users:   tt.fields.Users,
				Status:  tt.fields.Status,
			}
			if got := gs.CanUserVote(tt.args.userID); got != tt.want {
				t.Errorf("CanUserVote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSession_ApplyVote(t *testing.T) {
	type fields struct {
		ID      primitive.ObjectID
		OwnerId string
		Users   []*UserInfo
		Status  string
	}
	type args struct {
		votingUserID string
		votedUserID  string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          bool
		expectedUsers []*UserInfo
	}{
		{
			name: "Dani votes tomi",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId:   "dani",
						Alive:    true,
						Role:     ROLE_MAFIA,
						HasVoted: false,
					},
					{
						UserId:   "tomi",
						Alive:    true,
						Role:     ROLE_CITIZEN,
						HasVoted: false,
					},
				},
				Status: STAGE_DISCUSSION,
			},
			args: args{
				votingUserID: "dani",
				votedUserID:  "tomi",
			},
			want: true,
			expectedUsers: []*UserInfo{
				{
					UserId:   "dani",
					Alive:    true,
					Role:     ROLE_MAFIA,
					Votes:    0,
					HasVoted: true,
				},
				{
					UserId:   "tomi",
					Alive:    true,
					Role:     ROLE_CITIZEN,
					HasVoted: false,
					Votes:    1,
				},
			},
		},
		{
			name: "Dani votes inexistent username",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId:   "dani",
						Alive:    true,
						Role:     ROLE_MAFIA,
						HasVoted: false,
					},
					{
						UserId:   "tomi",
						Alive:    true,
						Role:     ROLE_CITIZEN,
						HasVoted: false,
					},
				},
				Status: STAGE_DISCUSSION,
			},
			args: args{
				votingUserID: "dani",
				votedUserID:  "invalid-user",
			},
			want: false,
			expectedUsers: []*UserInfo{
				{
					UserId:   "dani",
					Alive:    true,
					Role:     ROLE_MAFIA,
					HasVoted: true,
				},
				{
					UserId:   "tomi",
					Alive:    true,
					Role:     ROLE_CITIZEN,
					HasVoted: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GameSession{
				ID:      tt.fields.ID,
				OwnerId: tt.fields.OwnerId,
				Users:   tt.fields.Users,
				Status:  tt.fields.Status,
			}
			if got := gs.ApplyVote(tt.args.votingUserID, tt.args.votedUserID); got != tt.want {
				t.Errorf("ApplyVote() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(gs.Users, tt.expectedUsers) {
				t.Errorf("gs.users = %v, want %v", gs.Users, tt.expectedUsers)
			}
		})
	}
}

func TestGameSession_CanUserStartTheGame(t *testing.T) {
	type fields struct {
		ID      primitive.ObjectID
		OwnerId string
		Users   []*UserInfo
		Status  string
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "User starts a game that hasn't begun yet",
			fields: fields{
				OwnerId: "mily",
				Users: []*UserInfo{
					{
						UserId: "mily",
					},
					{
						UserId: "tomi",
					},
				},
				Status: STAGE_PENDING,
			},
			want: true,
		},
		{
			name: "User starts a game that's already begun",
			fields: fields{
				OwnerId: "tomi",
				Users: []*UserInfo{
					{
						UserId: "mily",
					},
					{
						UserId: "tomi",
					},
				},
				Status: STAGE_DISCUSSION,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameSession{
				ID:      tt.fields.ID,
				OwnerId: tt.fields.OwnerId,
				Users:   tt.fields.Users,
				Status:  tt.fields.Status,
			}
			if got := gs.CanUserStartTheGame(); got != tt.want {
				t.Errorf("CanUserStartTheGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSession_StartGame(t *testing.T) {
	type fields struct {
		ID      primitive.ObjectID
		OwnerId string
		Users   []*UserInfo
		Status  string
	}
	type want struct {
		output        bool
		status        string
		citizenAmount int
		mafiaAmount   int
		policeAmount  int
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Not enough users to start the game",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani",
					},
					{
						UserId: "tomi",
					},
				},
				Status: STAGE_PENDING,
			},
			want: want{
				output:        false,
				status:        STAGE_PENDING,
				citizenAmount: 0,
				mafiaAmount:   0,
				policeAmount:  0,
			},
		},
		{
			name: "Min users to start the game",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani",
					},
					{
						UserId: "tomi",
					},
					{
						UserId: "mili",
					},
				},
				Status: STAGE_PENDING,
			},
			want: want{
				output:        true,
				status:        STAGE_MAFIA,
				citizenAmount: 1,
				mafiaAmount:   1,
				policeAmount:  1,
			},
		},
		{
			name: "Nine users to start the game",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani1",
					},
					{
						UserId: "tomi2",
					},
					{
						UserId: "mili3",
					},
					{
						UserId: "dani4",
					},
					{
						UserId: "tomi5",
					},
					{
						UserId: "mili6",
					},
					{
						UserId: "dani7",
					},
					{
						UserId: "tomi8",
					},
					{
						UserId: "mili9",
					},
				},
				Status: STAGE_PENDING,
			},
			want: want{
				output:        true,
				status:        STAGE_MAFIA,
				citizenAmount: 5,
				mafiaAmount:   2,
				policeAmount:  2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GameSession{
				ID:      tt.fields.ID,
				OwnerId: tt.fields.OwnerId,
				Users:   tt.fields.Users,
				Status:  tt.fields.Status,
			}
			got := gs.StartGame()
			if got != tt.want.output {
				t.Errorf("StartGame() = %v, want %v", got, tt.want.output)
			}
			if gs.Status != tt.want.status {
				t.Errorf("Game status = %v, want %v", gs.Status, tt.want.status)
			}
			citizenAmountGot := 0
			policeAmountGot := 0
			mafiaAmountGot := 0

			for _, user := range gs.Users {
				switch user.Role {
				case ROLE_POLICE:
					policeAmountGot++
				case ROLE_MAFIA:
					mafiaAmountGot++
				case ROLE_CITIZEN:
					citizenAmountGot++
				}
			}
			if citizenAmountGot != tt.want.citizenAmount || policeAmountGot != tt.want.policeAmount ||
				mafiaAmountGot != tt.want.mafiaAmount {
				t.Errorf("Users roles = citizen:%v, mafia:%v, police:%v. Want citizen:%v, mafia:%v, police:%v.",
					citizenAmountGot, mafiaAmountGot, policeAmountGot, tt.want.citizenAmount, tt.want.mafiaAmount, tt.want.policeAmount)
			}
		})
	}
}
