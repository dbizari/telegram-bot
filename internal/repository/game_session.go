package repository

import (
	"context"
	"github.com/pkg/errors"
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
