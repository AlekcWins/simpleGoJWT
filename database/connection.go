package database

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"simpleGoJWT/utils"
)

var DB *mongo.Database

func NewClient(ctx context.Context, host, port, username, password, database, authDb string) (db *mongo.Database, err error) {
	var mongoDbURL string
	if username == "" && password == "" {
		mongoDbURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		mongoDbURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}
	utils.GetLogger().Info(mongoDbURL)
	clientOptions := options.Client().ApplyURI(mongoDbURL)

	if username != "" && password != "" {
		if authDb == "" {
			authDb = database
		}

		clientOptions.SetAuth(options.Credential{
			AuthSource: authDb,
			Username:   username,
			Password:   password,
		})
	}
	//Connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.New("failed connect mongo db")
	}

	//Ping
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to ping mongo db")
	}
	return client.Database(database), nil
}

func Connect(host, port, username, password, database, authDb string) {
	connect, err := NewClient(context.Background(), host, port, username, password, database, authDb)
	if err != nil {
		utils.GetLogger().Fatal("cant connect mongo db")
	}
	DB = connect
}
