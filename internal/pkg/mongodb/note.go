package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"yap-pwkeeper/internal/app/server/wallet"
	"yap-pwkeeper/internal/pkg/models"
)

func (db *Mongodb) AddNote(ctx context.Context, note models.Note) (string, error) {
	coll := db.client.Database(dbName).Collection(collNotes)
	res, err := coll.InsertOne(ctx, note)
	if err != nil {
		return "", err
	}
	oid, err := oid2string(res.InsertedID)
	return oid, err
}

func (db *Mongodb) GetNote(ctx context.Context, docId string, userId string) (models.Note, error) {
	coll := db.client.Database(dbName).Collection(collNotes)
	note := models.Note{}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return note, err
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", userId},
	}
	if err := coll.FindOne(ctx, filter).Decode(&note); err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			err = wallet.ErrNotFound
		}
		return note, err
	}
	//note, ok := result.(models.Note)
	//if !ok {
	//	err = ErrBadDoc
	//}
	return note, err
}

func (db *Mongodb) ModifyNote(ctx context.Context, note models.Note) error {
	coll := db.client.Database(dbName).Collection(collNotes)
	id, err := primitive.ObjectIDFromHex(note.Id)
	if err != nil {
		return err
	}
	newNote := struct {
		Id          primitive.ObjectID `bson:"_id"`
		models.Note `bson:"inline"`
	}{
		Id:   id,
		Note: note,
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", note.UserId},
	}
	var result interface{}
	err = coll.FindOneAndReplace(ctx, filter, newNote).Decode(&result)
	if errors.Is(mongo.ErrNoDocuments, err) {
		err = wallet.ErrNotFound
	}
	return err
}
