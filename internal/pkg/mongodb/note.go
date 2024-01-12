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

// AddNote places new Note in the database
func (db *Mongodb) AddNote(ctx context.Context, note models.Note) (string, error) {
	coll := db.client.Database(dbName).Collection(collNotes)
	res, err := coll.InsertOne(ctx, note)
	if err != nil {
		return "", err
	}
	oid, err := oid2string(res.InsertedID)
	return oid, err
}

// GetNote returns Note from database
func (db *Mongodb) GetNote(ctx context.Context, docId string, userId string) (models.Note, error) {
	coll := db.client.Database(dbName).Collection(collNotes)
	note := models.Note{}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return note, err
	}
	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "user_id", Value: userId},
	}
	if err := coll.FindOne(ctx, filter).Decode(&note); err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			err = documents.ErrNotFound
		}
		return note, err
	}
	return note, err
}

// ModifyNote updates record in database. Also called in delete action, because deleted
// documents are only marked for with a flag, but not actually deleted.
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
		{Key: "_id", Value: id},
		{Key: "user_id", Value: note.UserId},
	}
	var result interface{}
	err = coll.FindOneAndReplace(ctx, filter, newNote).Decode(&result)
	if errors.Is(mongo.ErrNoDocuments, err) {
		err = documents.ErrNotFound
	}
	return err
}

// GetNotesStream produces stream of Notes updates, happened between minSerial and maxSerial
func (db *Mongodb) GetNotesStream(ctx context.Context, userId string, minSerial, maxSerial int64, chData chan interface{}) error {
	coll := db.client.Database(dbName).Collection(collNotes)
	filter := bson.D{
		{Key: "user_id", Value: userId},
		{Key: "serial", Value: bson.D{{Key: "$gt", Value: minSerial}}},
		{Key: "serial", Value: bson.D{{Key: "$lt", Value: maxSerial}}},
	}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer func() { _ = cursor.Close(context.Background()) }()
	for cursor.Next(ctx) {
		var note models.Note
		if err := cursor.Decode(&note); err != nil {
			return err
		}
		chData <- note
	}
	return nil
}
