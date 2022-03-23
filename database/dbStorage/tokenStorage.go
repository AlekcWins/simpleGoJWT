package dbStorage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"simpleGoJWT/database/dbModel"
	"simpleGoJWT/interfaces/database"
)

type TokenStorage struct {
	collection *mongo.Collection
	ctx        context.Context
}

func (ts *TokenStorage) CreateOrUpdate(token dbModel.Token) (t dbModel.Token, err error) {
	tokenFound, err := ts.FindOneWithUserID(token.UserUId)
	if err != nil {
		id, err := ts.Create(token)
		if err != nil {
			return t, errors.New("error create token field in mongo")
		}
		token.ID = id
		return token, nil
	}
	token.ID = tokenFound.ID
	err = ts.Update(token)
	if err != nil {
		return t, errors.New("error update token field in mongo")
	}
	return token, nil

}

func (ts *TokenStorage) Create(token dbModel.Token) (string, error) {
	result, err := ts.collection.InsertOne(ts.ctx, token)
	if err != nil {
		return "", fmt.Errorf("failed to create object: %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to convert objectid to hex. probably oid: %s", oid)
}

func (ts *TokenStorage) FindOne(id string) (t dbModel.Token, err error) {
	filter := bson.M{"_id": id}

	result := ts.collection.FindOne(ts.ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return t, errors.New("object not found")
		}
		return t, fmt.Errorf("failed to find one token by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&t); err != nil {
		return t, fmt.Errorf("failed to decode token (id:%s) from DB due to error: %v", id, err)
	}

	return t, nil
}
func (ts *TokenStorage) FindOneWithUserID(userUId string) (t dbModel.Token, err error) {
	filter := bson.M{"user_id": userUId}

	result := ts.collection.FindOne(ts.ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return t, errors.New("object not found")
		}
		return t, fmt.Errorf("failed to find one token by id: %s due to error: %v", userUId, err)
	}
	if err = result.Decode(&t); err != nil {
		return t, fmt.Errorf("failed to decode token (id:%s) from DB due to error: %v", userUId, err)
	}

	return t, nil
}

func (ts *TokenStorage) Update(token dbModel.Token) error {
	objectID, err := primitive.ObjectIDFromHex(token.ID)
	if err != nil {
		return fmt.Errorf("failed to convert token ID to ObjectID. ID=%s", token.ID)
	}

	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marhsal token. error: %v", err)
	}

	var updateTokenObj bson.M
	err = bson.Unmarshal(userBytes, &updateTokenObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal token bytes. error: %v", err)
	}

	delete(updateTokenObj, "_id")

	update := bson.M{
		"$set": updateTokenObj,
	}

	result, err := ts.collection.UpdateOne(ts.ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update token query. error: %v", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("object not found failed to update")
	}

	return nil
}

func (ts *TokenStorage) Delete(userUId string) error {
	token, err := ts.FindOne(userUId)
	if err != nil {
		return fmt.Errorf("failed to find  token with user ID ID=%s", userUId)
	}
	objectID, err := primitive.ObjectIDFromHex(token.ID)
	if err != nil {
		return fmt.Errorf("failed to convert token ID to ObjectID. ID=%s", token.ID)
	}

	filter := bson.M{"_id": objectID}

	result, err := ts.collection.DeleteOne(ts.ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("object not found failed to delete")
	}

	return nil
}

func NewTokenStorage(database *mongo.Database) database.TokenStorage {
	return &TokenStorage{
		collection: database.Collection("tokens"),
		ctx:        context.TODO(),
	}
}
