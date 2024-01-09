package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"yap-pwkeeper/internal/app/server/documents"
	"yap-pwkeeper/internal/pkg/models"
)

func (db *Mongodb) AddCredential(ctx context.Context, credential models.Credential) (string, error) {
	coll := db.client.Database(dbName).Collection(collCredentials)
	res, err := coll.InsertOne(ctx, credential)
	if err != nil {
		return "", err
	}
	oid, err := oid2string(res.InsertedID)
	return oid, err
}

func (db *Mongodb) GetCredential(ctx context.Context, docId string, userId string) (models.Credential, error) {
	coll := db.client.Database(dbName).Collection(collCredentials)
	credential := models.Credential{}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return credential, err
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", userId},
	}
	if err := coll.FindOne(ctx, filter).Decode(&credential); err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			err = documents.ErrNotFound
		}
		return credential, err
	}
	return credential, err
}

func (db *Mongodb) ModifyCredential(ctx context.Context, credential models.Credential) error {
	coll := db.client.Database(dbName).Collection(collCredentials)
	id, err := primitive.ObjectIDFromHex(credential.Id)
	if err != nil {
		return err
	}
	newCredential := struct {
		Id                primitive.ObjectID `bson:"_id"`
		models.Credential `bson:"inline"`
	}{
		Id:         id,
		Credential: credential,
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", credential.UserId},
	}
	var result interface{}
	err = coll.FindOneAndReplace(ctx, filter, newCredential).Decode(&result)
	if errors.Is(mongo.ErrNoDocuments, err) {
		err = documents.ErrNotFound
	}
	return err
}

func (db *Mongodb) GetCredentialsStream(ctx context.Context, userId string, minSerial, maxSerial int64, chData chan interface{}) error {
	coll := db.client.Database(dbName).Collection(collCredentials)
	filter := bson.D{
		{"user_id", userId},
		{"serial", bson.D{{"$gt", minSerial}}},
		{"serial", bson.D{{"$lt", maxSerial}}},
	}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer func() { _ = cursor.Close(context.Background()) }()
	for cursor.Next(ctx) {
		var cred models.Credential
		if err := cursor.Decode(&cred); err != nil {
			return err
		}
		chData <- cred
	}
	return nil
}
