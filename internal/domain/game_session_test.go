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

func TestGameSession_GetRole(t *testing.T) {
	type fields struct {
		ID      primitive.ObjectID
		OwnerId string
		Users   []*UserInfo
		Status  string
	}
	type args struct {
		userId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "User gets their role",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani",
						Role:   ROLE_MAFIA,
					},
					{
						UserId: "tomi",
						Role:   ROLE_CITIZEN,
					},
				},
			},
			args: args{userId: "dani"},
			want: ROLE_MAFIA,
		},
		{
			name: "User doesn't have a role",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani",
					},
				},
			},
			args: args{userId: "dani"},
			want: "",
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
			if got := gs.GetRole(tt.args.userId); got != tt.want {
				t.Errorf("GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSession_CanUserAskForRole(t *testing.T) {
	type fields struct {
		ID      primitive.ObjectID
		OwnerId string
		Users   []*UserInfo
		Status  string
	}
	type args struct {
		userId    string
		userToAsk string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Police asks for role on their turn",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani",
						Role:   ROLE_POLICE,
					},
					{
						UserId: "mili",
						Role:   ROLE_MAFIA,
					},
				},
				Status: STAGE_POLICE,
			},
			args: args{userId: "dani", userToAsk: "mili"},
			want: true,
		},
		{
			name: "Police asks for role when it's not their turn",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani",
						Role:   ROLE_POLICE,
					},
					{
						UserId: "mili",
						Role:   ROLE_MAFIA,
					},
				},
				Status: STAGE_DISCUSSION,
			},
			args: args{userId: "dani", userToAsk: "mili"},
			want: false,
		},
		{
			name: "User asks for their own role",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani",
						Role:   ROLE_POLICE,
					},
					{
						UserId: "mili",
						Role:   ROLE_MAFIA,
					},
				},
				Status: STAGE_DISCUSSION,
			},
			args: args{userId: "mili", userToAsk: "mili"},
			want: true,
		},
		{
			name: "User asks for another user's role",
			fields: fields{
				Users: []*UserInfo{
					{
						UserId: "dani",
						Role:   ROLE_CITIZEN,
					},
					{
						UserId: "mili",
						Role:   ROLE_MAFIA,
					},
				},
				Status: STAGE_DISCUSSION,
			},
			args: args{userId: "dani", userToAsk: "mili"},
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
			if got := gs.CanUserAskForRole(tt.args.userId, tt.args.userToAsk); got != tt.want {
				t.Errorf("CanUserAskForRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
