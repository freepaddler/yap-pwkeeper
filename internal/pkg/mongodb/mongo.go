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

	"yap-pwkeeper/internal/pkg/aaa"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

const (
	connTimeout time.Duration = 5 * time.Second
	dbName                    = "pwkeeper"
	collUsers                 = "aaa"
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

func (db *Mongodb) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}

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
			{"login", login},
			{"state", models.StateActive},
		}).Decode(&user)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			err = aaa.ErrNotFound
		}
	}
	return user, err
}

// createIndexes creates database indexes
func (db *Mongodb) createIndexes(ctx context.Context) error {
	// unique index login@Users
	coll := db.client.Database(dbName).Collection(collUsers)
	userLogin := mongo.IndexModel{
		Keys:    bson.D{{"login", 1}},
		Options: options.Index().SetUnique(true),
	}
	logger.Log().Infof("creating index login_1_unique on %s", collUsers)
	_, err := coll.Indexes().CreateOne(ctx, userLogin)
	if err != nil {
		return err
	}
	return nil
}

func oid2string(oid interface{}) (string, error) {
	if id, ok := oid.(primitive.ObjectID); !ok {
		return "", errors.New("failed to get id from oid")
	} else {
		return id.Hex(), nil
	}
}
