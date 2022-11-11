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

func (gsr gameSessionRepository) AddPlayer(ctx context.Context, sessionId string, newUser *domain.UserInfo) (string, error) {

	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(sessionId)
	filter := bson.D{{"_id", id}}
	var session domain.GameSession
	err = collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return "", err
	}

	if session.ID.IsZero() {
		return "", errors.New("sesion not found")
	}

	for _, user := range session.Users {
		if user.UserId == newUser.UserId {
			return "", errors.New("¡You are already in this game!")
		}
	}

	update := bson.D{{"$set", bson.D{{"users", append(session.Users, *newUser)}}}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return session.ID.Hex(), nil
}

func (gsr gameSessionRepository) ExitGame(ctx context.Context, userName string) (bool, error) {

	collection := gsr.Database("mafia").Collection("game_session")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.D{{"users", bson.D{{"$elemMatch", bson.D{{
		"user_id",
		"tfanciotti"}}}}}}
	var session domain.GameSession
	err := collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return false, nil
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
		return false, err
	}

	return true, nil
}
