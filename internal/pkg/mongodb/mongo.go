package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"yap-pwkeeper/internal/pkg/logger"
)

const (
	connTimeout     time.Duration = 5 * time.Second
	dbName                        = "pwkeeper"
	collUsers                     = "users"
	collSerials                   = "serial"
	collNotes                     = "notes"
	collCredentials               = "credentials"
	collCards                     = "cards"
	collFiles                     = "files"
)

var (
	ErrBadId = errors.New("invalid documentId")
	//ErrBadDoc = errors.New("requested and returned document models do not match")
)

type Mongodb struct {
	uri    string
	client *mongo.Client
}

func New(ctx context.Context, uri string, opts ...func(db *Mongodb)) (*Mongodb, error) {
	db := new(Mongodb)
	var err error
	db.uri = uri
	if db.client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri).SetSocketTimeout(connTimeout)); err != nil {
		return nil, err
	}
	if err := db.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}
	return db, nil
}

// Close closes DB connection
func (db *Mongodb) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

// createIndexes creates database indexes
func (db *Mongodb) createIndexes(ctx context.Context) error {
	// unique index login @users
	coll := db.client.Database(dbName).Collection(collUsers)
	userLogin := mongo.IndexModel{
		Keys:    bson.D{{"login", 1}},
		Options: options.Index().SetUnique(true),
	}
	logger.Log().Infof("create index: login 1 unique for collection %s", collUsers)
	_, err := coll.Indexes().CreateOne(ctx, userLogin)
	if err != nil {
		return err
	}

	// documents search indexes
	search := mongo.IndexModel{
		Keys: bson.D{{"serial", -1}, {"user_id", 1}},
	}
	for _, v := range []string{collCards, collNotes, collSerials} {
		coll = db.client.Database(dbName).Collection(v)
		logger.Log().Infof("create index: serial -1 user_id 1 for collection %s", v)
		_, err := coll.Indexes().CreateOne(ctx, search)
		if err != nil {
			return err
		}
	}

	return nil
}

func oid2string(oid interface{}) (string, error) {
	if id, ok := oid.(primitive.ObjectID); !ok {
		return "", ErrBadId
	} else {
		return id.Hex(), nil
	}
}
