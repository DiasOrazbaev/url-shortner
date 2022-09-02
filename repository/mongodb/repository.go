package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/DiasOrazbaev/url-shortner/shortner"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoUrl string, mongoTimeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoUrl, mongoDB string, mongoTimeout int) (shortner.RedirectRepository, error) {
	repo := &mongoRepository{
		database: mongoDB,
		timeout:  time.Duration(mongoTimeout) * time.Second,
	}
	client, err := newMongoClient(mongoUrl, time.Duration(mongoTimeout)*time.Second)
	if err != nil {
		return nil, errors.New(err.Error() + " repository.mongodb.NewMongoRepository")
	}
	repo.client = client
	return repo, nil
}

func (m *mongoRepository) Find(code string) (*shortner.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	redirect := &shortner.Redirect{}
	collection := m.client.Database(m.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(shortner.ErrRedirectNotFound.Error() + " repository.mongodb.mongoRepository.Find")
		}
		return nil, errors.New(err.Error() + " repository.mongodb.mongoRepository.Find")
	}
	return redirect, nil
}

func (m *mongoRepository) Save(redirect *shortner.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()
	collection := m.client.Database(m.database).Collection("redirects")
	_, err := collection.InsertOne(ctx, bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	})
	if err != nil {
		return errors.New(err.Error() + " repository.mongodb.mongoRepository.Save")
	}
	return nil
}
