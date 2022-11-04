package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	client ClientAPI
	once   sync.Once
)

type ClientAPI interface {
	CreateGame()
}

type gameSessionRepository struct {
	*mongo.Client
}

func GetGameSessionRepositoryClient() ClientAPI {
	once.Do(func() {
		mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			// ToDo handle error
			panic(err)
		}

		client = &gameSessionRepository{
			Client: mongoClient,
		}

		// ToDo add all the necessary things
	})

	return client
}

func (gsr gameSessionRepository) CreateGame() {

}
