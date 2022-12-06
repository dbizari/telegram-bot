package repository

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"tdl/internal/domain"
	"time"
)

var (
	client GameSessionRepositoryAPI
	once   sync.Once
)

type GameSessionRepositoryAPI interface {
	CreateGame(ctx context.Context, gameSession *domain.GameSession) (string, error)
	Get(ctx context.Context, gameSessionID string) (*domain.GameSession, error)
	Update(ctx context.Context, gameSession *domain.GameSession) error
	GetByMember(ctx context.Context, username string) (*domain.GameSession, error)
	FindGame(ctx context.Context, sessionId string) (domain.GameSession, error)
	AddPlayer(ctx context.Context, sessionId string, userInfo *domain.UserInfo) (string, error)
	ExitGame(ctx context.Context, userName string) (bool, error)
}

type gameSessionRepository struct {
	*mongo.Client
}

func GetGameSessionRepositoryClient() GameSessionRepositoryAPI {
	once.Do(func() {
		ctx := context.Background()
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://fiuba:eldanisape@mafia-bot.dx3pk6a.mongodb.net/?retryWrites=true&w=majority")) // Todo change this for env vars
		if err != nil {
			err = errors.Wrap(err, "failed to create mongodb client")
			panic(err)
		}

		err = mongoClient.Ping(ctx, readpref.Primary())
		if err != nil {
			err = errors.Wrap(err, "ping failed to mongo")
			panic(err)
		}

		client = &gameSessionRepository{
			Client: mongoClient,
		}
	})

	return client
}

func (gsr gameSessionRepository) CreateGame(ctx context.Context, gameSession *domain.GameSession) (string, error) {
	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, gameSession)
	if err != nil {
		err = errors.Wrap(err, "error trying to create game")
		return "", err
	}

	id, _ := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (gsr gameSessionRepository) Get(ctx context.Context, gameSessionID string) (*domain.GameSession, error) {
	id, err := primitive.ObjectIDFromHex(gameSessionID)
	if err != nil {
		return nil, errors.Wrap(err, "error trying to convert gameSessionID to ObjectID")
	}

	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var session domain.GameSession
	filter := bson.D{{"_id", id}}
	err = collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error trying to get game session from db")
	}

	return &session, nil
}

func (gsr gameSessionRepository) FindGame(ctx context.Context, sessionId string) (domain.GameSession, error) {

	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	session := domain.GameSession{}

	// Find game
	id, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		return session, errors.Wrap(err, "error trying to handle sessionId")
	}

	filter := bson.D{{"_id", id}}
	err = collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return session, errors.New("Game not found")
		}
		return session, errors.Wrap(err, "error trying to get game session from db")
	}

	return session, nil
}

func (gsr gameSessionRepository) GetByMember(ctx context.Context, userID string) (*domain.GameSession, error) {
	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var session domain.GameSession
	filter := bson.D{
		{"users", bson.D{
			{"$elemMatch", bson.D{
				{"user_id", userID},
			}},
		}},
	}

	err := collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error trying to get game session from db")
	}

	return &session, nil
}

func (gsr gameSessionRepository) AddPlayer(ctx context.Context, sessionId string, newUser *domain.UserInfo) (string, error) {

	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if already player is in a game
	session := domain.GameSession{}
	filter := bson.D{{"users", bson.D{{"$elemMatch", bson.D{{
		"user_id",
		newUser.UserId}}}}}}
	err := collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return "", errors.Wrap(err, "error trying to find user on db")
		}
	}

	if !session.ID.IsZero() {
		return "", errors.New("You are already in a game!")
	}

	for _, user := range session.Users {
		if user.UserId == newUser.UserId {
			return "", errors.New("you are already in a game")
		}
	}

	// Find game
	session, err = gsr.FindGame(ctx, sessionId)
	if err != nil {
		return "", err
	}

	// add player to the game
	update := bson.D{{"$set", bson.D{{"users", append(session.Users, *newUser)}}}}
	_, err = collection.UpdateOne(ctx, session, update)
	if err != nil {
		return "", errors.Wrap(err, "error trying to update the game on db")
	}

	return session.ID.Hex(), nil
}

func (gsr gameSessionRepository) Update(ctx context.Context, gameSession *domain.GameSession) error {

	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": gameSession.ID}
	result, err := collection.ReplaceOne(ctx, filter, gameSession)
	if err != nil {
		return errors.Wrap(err, "error on update game session from db")
	}

	if result.ModifiedCount != 1 {
		return errors.New("error on update game session from db, the document was not updated")
	}

	return nil
}

func (gsr gameSessionRepository) ExitGame(ctx context.Context, userName string) (bool, error) {

	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.D{{"users", bson.D{{"$elemMatch", bson.D{{
		"user_id",
		userName}}}}}}
	var session domain.GameSession
	err := collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, errors.Wrap(err, "error trying to get game session from db")
	}

	update := bson.D{{
		"$pull",
		bson.D{{
			"users",
			bson.D{{
				"user_id",
				userName}}}}}}

	err = collection.FindOneAndUpdate(ctx, filter, update).Decode(&session)
	if err != nil {
		return false, errors.Wrap(err, "error trying to update game session on db")
	}

	return true, nil
}
