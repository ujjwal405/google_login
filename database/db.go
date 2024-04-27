package database

import (
	"context"
	"errors"
	"time"

	model "github.com/ujjwal405/google_login/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	errEmailexist    = errors.New("email already exists")
	errEmailnotexist = errors.New("email doesn't exists")
	Emailexist       = "email already exists"
	Emailnotexist    = "email doesn't exists"
)

type AllDatabase interface {
	DBCheckEmail(email string) error
	DBSignup(user model.UserSignup) error
	DBGetData(email string) (model.UserSignup, error)
	DBUpdate(userid string, newdata model.UserData) error
	DbData(userid string) (model.UserSignup, error)
}
type Database struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewDatabase(client *mongo.Client, collection *mongo.Collection) AllDatabase {
	return &Database{
		Client:     client,
		Collection: collection,
	}
}

func (db *Database) DBCheckEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	count, err := db.Collection.CountDocuments(ctx, bson.M{"email": email})
	defer cancel()
	if err != nil {
		return err
	}
	if count > 0 {
		err = errEmailexist
		return err
	}
	return errEmailnotexist
}

func (db *Database) DBSignup(user model.UserSignup) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	_, err := db.Collection.InsertOne(ctx, user)
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DBGetData(email string) (model.UserSignup, error) {
	var user model.UserSignup
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	err := db.DBCheckEmail(email)
	defer cancel()
	if err != nil {
		if err.Error() == Emailnotexist {
			return user, err
		} else if err.Error() != Emailexist && err.Error() != Emailnotexist {
			return user, err
		}
	}
	err = db.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (db *Database) DBUpdate(userid string, newdata model.UserData) error {
	var founduser model.UserSignup
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	err := db.Collection.FindOne(ctx, bson.M{"user_id": userid}).Decode(&founduser)
	defer cancel()
	if err != nil {
		return err
	}
	if founduser.Isvalid {
		newdata.Email = founduser.Email
	}
	filter := bson.M{"user_id": userid}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err = db.Collection.UpdateOne(ctx,
		filter,
		bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "email", Value: newdata.Email},
			{Key: "phone", Value: newdata.Phone},
			{Key: "user_name", Value: newdata.Username}}}},
		&opt,
	)

	if err != nil {
		return err
	}
	return nil
}
func (db *Database) DbData(userid string) (model.UserSignup, error) {
	var user model.UserSignup
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	err := db.Collection.FindOne(ctx, bson.M{"user_id": userid}).Decode(&user)
	defer cancel()

	if err != nil {
		return user, err
	}
	return user, nil
}
