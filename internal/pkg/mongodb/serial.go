package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetSerials returns next usable serial, reserves next n serials in db collection
func (db *Mongodb) GetSerials(ctx context.Context, n int) (int64, error) {
	coll := db.client.Database(dbName).Collection(collSerials)
	filter := bson.D{}
	update := bson.D{
		{"$inc", bson.D{{"next", int64(n)}}},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	var res bson.M
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&res)
	if errors.Is(mongo.ErrNoDocuments, err) {
		return 0, nil
	}
	return res["next"].(int64), err
}
