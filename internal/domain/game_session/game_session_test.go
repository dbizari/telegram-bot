package game_session

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"tdl/internal/domain/game_stages"
	user_pkg "tdl/internal/domain/user"
	"testing"
)

func TestGameSession_CanUserVote(t *testing.T) {
	type fields struct {
		ID      primitive.ObjectID
		OwnerId string
		Users   []*user_pkg.UserInfo
		Stage   game_stages.GameStage
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
				Users: []*user_pkg.UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     user_pkg.ROLE_MAFIA,
					HasVoted: false,
				}},
				Stage: game_stages.Mafia{},
			},
			args: args{
				userID: "dani",
			},
			want: true,
		},
		{
			name: "Dead user cannot vote",
			fields: fields{
				Users: []*user_pkg.UserInfo{{
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
				Users: []*user_pkg.UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     user_pkg.ROLE_CITIZEN,
					HasVoted: false,
				}},
				Stage: game_stages.Mafia{},
			},
			args: args{
				userID: "dani",
			},
			want: false,
		},
		{
			name: "Mafia person votes on Discussion stage",
			fields: fields{
				Users: []*user_pkg.UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     user_pkg.ROLE_MAFIA,
					HasVoted: false,
				}},
				Stage: game_stages.Discussion{},
			},
			args: args{
				userID: "dani",
			},
			want: true,
		},
		{
			name: "Citizen person votes on Discussion stage",
			fields: fields{
				Users: []*user_pkg.UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     user_pkg.ROLE_MAFIA,
					HasVoted: false,
				}},
				Stage: game_stages.Discussion{},
			},
			args: args{
				userID: "dani",
			},
			want: true,
		},
		{
			name: "Citizen person votes on Police stage",
			fields: fields{
				Users: []*user_pkg.UserInfo{{
					UserId:   "dani",
					Alive:    true,
					Role:     user_pkg.ROLE_MAFIA,
					HasVoted: false,
				}},
				Stage: game_stages.Police{},
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
				Stage:   tt.fields.Stage,
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
		Users   []*user_pkg.UserInfo
		Stage   game_stages.GameStage
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
		expectedUsers []*user_pkg.UserInfo
	}{
		{
			name: "Dani votes tomi",
			fields: fields{
				Users: []*user_pkg.UserInfo{
					{
						UserId:   "dani",
						Alive:    true,
						Role:     user_pkg.ROLE_MAFIA,
						HasVoted: false,
					},
					{
						UserId:   "tomi",
						Alive:    true,
						Role:     user_pkg.ROLE_CITIZEN,
						HasVoted: false,
					},
				},
				Stage: game_stages.Discussion{},
			},
			args: args{
				votingUserID: "dani",
				votedUserID:  "tomi",
			},
			want: true,
			expectedUsers: []*user_pkg.UserInfo{
				{
					UserId:   "dani",
					Alive:    true,
					Role:     user_pkg.ROLE_MAFIA,
					Votes:    0,
					HasVoted: true,
				},
				{
					UserId:   "tomi",
					Alive:    true,
					Role:     user_pkg.ROLE_CITIZEN,
					HasVoted: false,
					Votes:    1,
				},
			},
		},
		{
			name: "Dani votes inexistent username",
			fields: fields{
				Users: []*user_pkg.UserInfo{
					{
						UserId:   "dani",
						Alive:    true,
						Role:     user_pkg.ROLE_MAFIA,
						HasVoted: false,
					},
					{
						UserId:   "tomi",
						Alive:    true,
						Role:     user_pkg.ROLE_CITIZEN,
						HasVoted: false,
					},
				},
				Stage: game_stages.Discussion{},
			},
			args: args{
				votingUserID: "dani",
				votedUserID:  "invalid-user",
			},
			want: false,
			expectedUsers: []*user_pkg.UserInfo{
				{
					UserId:   "dani",
					Alive:    true,
					Role:     user_pkg.ROLE_MAFIA,
					HasVoted: true,
				},
				{
					UserId:   "tomi",
					Alive:    true,
					Role:     user_pkg.ROLE_CITIZEN,
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
				Stage:   tt.fields.Stage,
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
		Users   []*user_pkg.UserInfo
		Stage   game_stages.GameStage
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
				Users: []*user_pkg.UserInfo{
					{
						UserId: "dani",
						Role:   user_pkg.ROLE_MAFIA,
					},
					{
						UserId: "tomi",
						Role:   user_pkg.ROLE_CITIZEN,
					},
				},
			},
			args: args{userId: "dani"},
			want: user_pkg.ROLE_MAFIA,
		},
		{
			name: "User doesn't have a role",
			fields: fields{
				Users: []*user_pkg.UserInfo{
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
				Stage:   tt.fields.Stage,
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
		Users   []*user_pkg.UserInfo
		Stage   game_stages.GameStage
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
				Users: []*user_pkg.UserInfo{
					{
						UserId: "dani",
						Role:   user_pkg.ROLE_POLICE,
					},
					{
						UserId: "mili",
						Role:   user_pkg.ROLE_MAFIA,
					},
				},
				Stage: game_stages.Police{},
			},
			args: args{userId: "dani", userToAsk: "mili"},
			want: true,
		},
		{
			name: "Police asks for role when it's not their turn",
			fields: fields{
				Users: []*user_pkg.UserInfo{
					{
						UserId: "dani",
						Role:   user_pkg.ROLE_POLICE,
					},
					{
						UserId: "mili",
						Role:   user_pkg.ROLE_MAFIA,
					},
				},
				Stage: game_stages.Discussion{},
			},
			args: args{userId: "dani", userToAsk: "mili"},
			want: false,
		},
		{
			name: "User asks for another user's role",
			fields: fields{
				Users: []*user_pkg.UserInfo{
					{
						UserId: "dani",
						Role:   user_pkg.ROLE_CITIZEN,
					},
					{
						UserId: "mili",
						Role:   user_pkg.ROLE_MAFIA,
					},
				},
				Stage: game_stages.Discussion{},
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
				Stage:   tt.fields.Stage,
			}
			if got := gs.CanUserAskForRole(tt.args.userId); got != tt.want {
				t.Errorf("CanUserAskForRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSession_CanUserStartTheGame(t *testing.T) {
	type fields struct {
		ID      primitive.ObjectID
		OwnerId string
		Users   []*user_pkg.UserInfo
		Stage   game_stages.GameStage
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
				Users: []*user_pkg.UserInfo{
					{
						UserId: "mily",
					},
					{
						UserId: "tomi",
					},
				},
				Stage: game_stages.Pending{},
			},
			want: true,
		},
		{
			name: "User starts a game that's already begun",
			fields: fields{
				OwnerId: "tomi",
				Users: []*user_pkg.UserInfo{
					{
						UserId: "mily",
					},
					{
						UserId: "tomi",
					},
				},
				Stage: game_stages.Discussion{},
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
				Stage:   tt.fields.Stage,
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
		Users   []*user_pkg.UserInfo
		Stage   game_stages.GameStage
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
				Users: []*user_pkg.UserInfo{
					{
						UserId: "dani",
					},
					{
						UserId: "tomi",
					},
				},
				Stage: game_stages.Pending{},
			},
			want: want{
				output:        false,
				status:        game_stages.STAGE_PENDING,
				citizenAmount: 0,
				mafiaAmount:   0,
				policeAmount:  0,
			},
		},
		{
			name: "Min users to start the game",
			fields: fields{
				Users: []*user_pkg.UserInfo{
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
				Stage: game_stages.Pending{},
			},
			want: want{
				output:        true,
				status:        game_stages.STAGE_MAFIA,
				citizenAmount: 1,
				mafiaAmount:   1,
				policeAmount:  1,
			},
		},
		{
			name: "Nine users to start the game",
			fields: fields{
				Users: []*user_pkg.UserInfo{
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
				Stage: game_stages.Pending{},
			},
			want: want{
				output:        true,
				status:        game_stages.STAGE_MAFIA,
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
				Stage:   tt.fields.Stage,
			}
			got := gs.StartGame()
			if got != tt.want.output {
				t.Errorf("StartGame() = %v, want %v", got, tt.want.output)
			}
			if gs.Stage.GetStageName() != tt.want.status {
				t.Errorf("Game status = %v, want %v", gs.Stage, tt.want.status)
			}
			citizenAmountGot := 0
			policeAmountGot := 0
			mafiaAmountGot := 0

			for _, user := range gs.Users {
				switch user.Role {
				case user_pkg.ROLE_POLICE:
					policeAmountGot++
				case user_pkg.ROLE_MAFIA:
					mafiaAmountGot++
				case user_pkg.ROLE_CITIZEN:
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
