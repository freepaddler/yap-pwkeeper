package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"yap-pwkeeper/internal/app/server/aaa"
	"yap-pwkeeper/internal/pkg/models"
)

// AddUser creates new user in database, returns user with new id
func (db *Mongodb) AddUser(ctx context.Context, user models.User) (models.User, error) {
	coll := db.client.Database(dbName).Collection(collUsers)
	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return user, aaa.ErrDuplicate
		}
		return user, err
	}
	user.Id, err = oid2string(res.InsertedID)
	return user, err
}

// GetUserByLogin returns active user by specified login
func (db *Mongodb) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User
	coll := db.client.Database(dbName).Collection(collUsers)
	err := coll.FindOne(ctx,
		bson.D{
			{Key: "login", Value: login},
			{Key: "state", Value: models.StateActive},
		}).Decode(&user)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			err = aaa.ErrNotFound
		}
	}
	return user, err
}
