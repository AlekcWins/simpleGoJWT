package dbStorage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"simpleGoJWT/database/dbModel"
	"simpleGoJWT/interfaces/database"
)

type UserStorage struct {
	collection *mongo.Collection
	ctx        context.Context
}

func (us *UserStorage) FindOneWithEmail(email string) (u dbModel.User, err error) {

	filter := bson.M{"email": email}

	result := us.collection.FindOne(us.ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, errors.New("object not found")
		}
		return u, fmt.Errorf("failed to find one user by email: %s due to error: %v", email, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user (email:%s) from DB due to error: %v", email, err)
	}

	return u, nil
}

func (us *UserStorage) Create(user dbModel.User) (string, error) {
	result, err := us.collection.InsertOne(us.ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create object: %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to convert objectid to hex. probably oid: %s", oid)
}

func (us *UserStorage) FindAll() (u []dbModel.User, err error) {
	cursor, err := us.collection.Find(us.ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find all users due to error: %v", err)
	}

	if err = cursor.All(us.ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return u, nil
}

func (us *UserStorage) FindOne(id string) (u dbModel.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}

	filter := bson.M{"_id": oid}

	result := us.collection.FindOne(us.ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, errors.New("object not found")
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user (id:%s) from DB due to error: %v", id, err)
	}

	return u, nil
}

func (us *UserStorage) Update(user dbModel.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", user.ID)
	}

	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marhsal user. error: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes. error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := us.collection.UpdateOne(us.ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %v", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("object not found failed to update")
	}

	return nil

}

func (us *UserStorage) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", id)
	}

	filter := bson.M{"_id": objectID}

	result, err := us.collection.DeleteOne(us.ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}
	if result.DeletedCount == 0 {
		return errors.New("object not found failed to delete")
	}

	return nil
}

func NewUserStorage(database *mongo.Database) database.UserStorage {
	collection := database.Collection("users")
	_, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		panic("create index to users collection on email field")
		return nil
	}
	return &UserStorage{
		collection: collection,
		ctx:        context.TODO(),
	}
}
