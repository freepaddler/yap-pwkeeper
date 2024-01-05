package mongodb

import (
	"context"
	"time"

	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

func (db *Mongodb) AddNote(ctx context.Context, note models.NoteDocument) error {
	coll := db.client.Database(dbName).Collection(collNotes)
	note.CreatedAt = time.Now()
	note.ModifiedAt = time.Now()
	note.State = models.StateActive
	res, err := coll.InsertOne(ctx, note)
	if err != nil {
		return err
	}
	oid, err := oid2string(res.InsertedID)
	logger.Log().WithCtxRequestId(ctx).With("documentId", oid).Debug("note added")
	return err
}

//func (db *Mongodb) ReplaceNote(ctx context.Context, note models.Note) error {
//	coll := db.client.Database(dbName).Collection(collNotes)
//	id, err := primitive.ObjectIDFromHex(note.Id)
//	if err != nil {
//		logger.Log().WithErr(err).WithCtxRequestId(ctx).With("documentId", note.Id).Warn("note replace failed")
//		return err
//	}
//	replace := struct {
//		Id          primitive.ObjectID `bson:"_id"`
//		models.Note `bson:"inline"`
//	}{
//		Id:   id,
//		Note: note,
//	}
//	replace.ModifiedAt = time.Now()
//	updateFilter := bson.D{{"_id", id}, {"state", models.StateActive}}
//	res, err := coll.ReplaceOne(ctx, updateFilter, replace)
//	if err != nil {
//		logger.Log().WithErr(err).WithCtxRequestId(ctx).With("documentId", id.Hex()).Warn("note replace failed")
//		return err
//	}
//	if res.ModifiedCount == 1 {
//		logger.Log().WithCtxRequestId(ctx).With("documentId", id.Hex()).Debug("note replaced")
//	} else if res.MatchedCount > 0 {
//		logger.Log().WithCtxRequestId(ctx).With("documentId", id.Hex()).Warn("note replace?")
//	} else {
//		logger.Log().WithCtxRequestId(ctx).With("documentId", id.Hex()).Warn("note to replace not found")
//	}
//	return nil
//}
//
//func (db *Mongodb) DelNote(ctx context.Context, documentId string) error {
//	coll := db.client.Database(dbName).Collection(collNotes)
//	id, err := primitive.ObjectIDFromHex(documentId)
//	if err != nil {
//		logger.Log().WithErr(err).WithCtxRequestId(ctx).With("documentId", documentId).Warn("note delete failed")
//		return err
//	}
//	updateFilter := bson.D{{"_id", id}, {"state", models.StateActive}}
//	updateSet := bson.D{
//		{"$set", bson.D{
//			{"state", models.StateDeleted},
//			{"modified_at", time.Now()},
//		}},
//	}
//	res, err := coll.UpdateOne(ctx, updateFilter, updateSet)
//	if err != nil {
//		return err
//	}
//	if res.ModifiedCount == 1 {
//		logger.Log().WithCtxRequestId(ctx).With("documentId", documentId).Debug("note deleted")
//	} else if res.MatchedCount > 0 {
//		logger.Log().WithCtxRequestId(ctx).With("documentId", documentId).Warn("note already deleted")
//	} else {
//		logger.Log().WithCtxRequestId(ctx).With("documentId", documentId).Warn("note to delete not found")
//	}
//	return nil
//}
