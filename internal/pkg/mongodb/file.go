package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"yap-pwkeeper/internal/app/server/documents"
	"yap-pwkeeper/internal/pkg/models"
)

// AddFile stores file in database
func (db *Mongodb) AddFile(ctx context.Context, file models.File) (string, error) {
	coll := db.client.Database(dbName).Collection(collFiles)
	res, err := coll.InsertOne(ctx, file)
	if err != nil {
		return "", err
	}
	oid, err := oid2string(res.InsertedID)
	return oid, err
}

// ModifyFile replaces current document with new one. Used for update and for deletion
func (db *Mongodb) ModifyFile(ctx context.Context, file models.File) error {
	coll := db.client.Database(dbName).Collection(collFiles)
	id, err := primitive.ObjectIDFromHex(file.Id)
	if err != nil {
		return err
	}
	newFile := struct {
		Id          primitive.ObjectID `bson:"_id"`
		models.File `bson:"inline"`
	}{
		Id:   id,
		File: file,
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", file.UserId},
	}
	var result interface{}
	err = coll.FindOneAndReplace(ctx, filter, newFile).Decode(&result)
	if errors.Is(mongo.ErrNoDocuments, err) {
		err = documents.ErrNotFound
	}
	return err
}

// ModifyFileInfo updates record without touching file data and information.
// Should be used, when update was without file
func (db *Mongodb) ModifyFileInfo(ctx context.Context, file models.File) error {
	coll := db.client.Database(dbName).Collection(collFiles)
	id, err := primitive.ObjectIDFromHex(file.Id)
	if err != nil {
		return err
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", file.UserId},
	}
	update := bson.D{
		{"$set", bson.D{{"serial", file.Serial}}},
		{"$set", bson.D{{"name", file.Name}}},
		{"$set", bson.D{{"metadata", file.Metadata}}},
	}
	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return documents.ErrNotFound
	}
	return err
}

// GetFile returns file with binary data
func (db *Mongodb) GetFile(ctx context.Context, docId string, userId string) (models.File, error) {
	return db.getFile(ctx, docId, userId, false)
}

// GetFileInfo returns file without binary data
func (db *Mongodb) GetFileInfo(ctx context.Context, docId string, userId string) (models.File, error) {
	return db.getFile(ctx, docId, userId, true)
}

func (db *Mongodb) getFile(ctx context.Context, docId string, userId string, withoutData bool) (models.File, error) {
	coll := db.client.Database(dbName).Collection(collFiles)
	file := models.File{}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return file, err
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", userId},
	}
	opts := options.FindOne()
	if withoutData {
		opts = opts.SetProjection(bson.D{{"data", 0}})
	}
	if err := coll.FindOne(ctx, filter, opts).Decode(&file); err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			err = documents.ErrNotFound
		}
		return file, err
	}
	return file, err
}

// GetFilesInfoStream produces stream of FileInfo updates, happened between minSerial and maxSerial
func (db *Mongodb) GetFilesInfoStream(ctx context.Context, userId string, minSerial, maxSerial int64, chData chan interface{}) error {
	coll := db.client.Database(dbName).Collection(collFiles)
	filter := bson.D{
		{"user_id", userId},
		{"serial", bson.D{{"$gt", minSerial}}},
		{"serial", bson.D{{"$lt", maxSerial}}},
	}
	opts := options.Find().SetProjection(bson.D{{"data", 0}})
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer func() { _ = cursor.Close(context.Background()) }()
	for cursor.Next(ctx) {
		var file models.File
		if err := cursor.Decode(&file); err != nil {
			return err
		}
		chData <- file
	}
	return nil
}
