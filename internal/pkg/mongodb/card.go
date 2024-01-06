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

func (db *Mongodb) AddCard(ctx context.Context, card models.Card) (string, error) {
	coll := db.client.Database(dbName).Collection(collCards)
	res, err := coll.InsertOne(ctx, card)
	if err != nil {
		return "", err
	}
	oid, err := oid2string(res.InsertedID)
	return oid, err
}

func (db *Mongodb) GetCard(ctx context.Context, docId string, userId string) (models.Card, error) {
	coll := db.client.Database(dbName).Collection(collCards)
	card := models.Card{}
	id, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return card, err
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", userId},
	}
	if err := coll.FindOne(ctx, filter).Decode(&card); err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			err = wallet.ErrNotFound
		}
		return card, err
	}
	return card, err
}

func (db *Mongodb) ModifyCard(ctx context.Context, card models.Card) error {
	coll := db.client.Database(dbName).Collection(collCards)
	id, err := primitive.ObjectIDFromHex(card.Id)
	if err != nil {
		return err
	}
	newCard := struct {
		Id          primitive.ObjectID `bson:"_id"`
		models.Card `bson:"inline"`
	}{
		Id:   id,
		Card: card,
	}
	filter := bson.D{
		{"_id", id},
		{"user_id", card.UserId},
	}
	var result interface{}
	err = coll.FindOneAndReplace(ctx, filter, newCard).Decode(&result)
	if errors.Is(mongo.ErrNoDocuments, err) {
		err = wallet.ErrNotFound
	}
	return err
}

func (db *Mongodb) GetCardsStream(ctx context.Context, userId string, minSerial, maxSerial int64, chData chan interface{}) error {
	coll := db.client.Database(dbName).Collection(collCards)
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
		var card models.Card
		if err := cursor.Decode(&card); err != nil {
			return err
		}
		chData <- card
	}
	return nil
}
