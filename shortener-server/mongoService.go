package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	URL_UID = "_id"
)

type UrlDocument struct {
	Url string `bson:"Url"`
	Id  string `bson:"_id"`
}

type DbOptions struct {
	Database      string
	UrlCollection string
	UserName      string
	Password      string
	Host          string
	Port          int
}

//DbIface implementation for MongoDB
type MongoService struct {
	context     context.Context
	client      *mongo.Client
	collHandler *mongo.Collection
}

func NewMongoService(opt *DbOptions) (*MongoService, error) {
	//The port is optional so it is empty if none is provided to avoid generating a wrong MongoDBUri
	port := ""
	if opt.Port > 0 {
		port = fmt.Sprintf(":%v", opt.Port)
	}
	mongoDbUri := fmt.Sprintf("mongodb+srv://%s:%s@%s%s/?retryWrites=true&w=majority",
		opt.UserName,
		opt.Password,
		opt.Host,
		port)
	context := context.TODO()
	client, err := mongo.Connect(
		context,
		options.Client().ApplyURI(mongoDbUri).SetServerSelectionTimeout(20*time.Second))
	if err != nil {
		return nil, fmt.Errorf("error while creating mongoDB client: %w", err)
	}
	if err := client.Ping(context, nil); err != nil {
		return nil, fmt.Errorf("error while doing ping: %w", err)
	}
	return &MongoService{
		context:     context,
		client:      client,
		collHandler: client.Database(opt.Database).Collection(opt.UrlCollection),
	}, nil
}

func (ms *MongoService) ShutDown() error {
	return ms.client.Disconnect(ms.context)
}

func (ms *MongoService) GetItem(key string) (interface{}, error) {
	var urlDoc UrlDocument
	err := ms.collHandler.FindOne(
		ms.context,
		bson.D{primitive.E{Key: URL_UID, Value: key}}).Decode(&urlDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error while finding item: %w", err)
	}
	return urlDoc.Url, nil
}

func (ms *MongoService) SetItem(key string, value interface{}) error {
	newUrl := UrlDocument{
		Id:  key,
		Url: fmt.Sprintf("%v", value),
	}
	filter := bson.D{{Key: URL_UID, Value: key}}
	update := bson.M{"$set": newUrl}
	options := options.Update().SetUpsert(true)
	if _, err := ms.collHandler.UpdateOne(ms.context, filter, update, options); err != nil {
		return fmt.Errorf("error while setting item: %w", err)
	}
	return nil
}

func (ms *MongoService) Exists(key string) bool {
	if result := ms.collHandler.FindOne(ms.context, bson.D{{Key: URL_UID, Value: key}}); result.Err() != nil {
		return false
	}
	return true
}
