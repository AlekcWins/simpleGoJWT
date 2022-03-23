package services

import (
	"go.mongodb.org/mongo-driver/mongo"
	"simpleGoJWT/config"
	db "simpleGoJWT/database"
	"simpleGoJWT/database/dbStorage"
	"simpleGoJWT/interfaces/database"
	"sync"
)

func init() {
	cfg := config.GetConfig()
	dbConfig := cfg.MongoDB
	db.Connect(dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.AuthDB)
}

type dbServices struct {
	TokenStorage database.TokenStorage
	UserStorage  database.UserStorage
}

func newDbServices(database *mongo.Database) *dbServices {
	return &dbServices{
		TokenStorage: dbStorage.NewTokenStorage(database),
		UserStorage:  dbStorage.NewUserStorage(database),
	}
}

var instance *dbServices
var once sync.Once

func GetDbServices(database *mongo.Database) *dbServices {
	once.Do(func() {
		instance = newDbServices(database)
	})
	return instance
}
