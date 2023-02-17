package datastore

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
)

var (
	Db               *MongoDatabase
)

type MongoDatabase struct {
	Url    string
	Name   string
	Client *mongo.Client
	Jobs   string
}

func GetDatastore() MongoDB {
	return Db
}

func InitialiseAndConnectToMongo(url, username, password, dbName string) *MongoDatabase {
	var md MongoDatabase
	md = md.Bootstrap(url, username, password, dbName).(MongoDatabase)

	md.Connect()
	Db = &md
	return &md
}

func (m MongoDatabase) Connect() {
	err := m.Client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}
	err = m.Client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	log.Println("Connected to the Database")
}

func (m MongoDatabase) Bootstrap(url, username, password, dbName string) interface{} {
	finalUrl := "mongodb+srv://" + username + ":" + password + "@" + url

	log.Println("Connecting to mongo")

	client, err := mongo.NewClient(options.Client().ApplyURI(finalUrl), options.Client().SetMaxPoolSize(5))
	if err != nil {
		panic(err)
	}
	m.Name = dbName
	m.Client = client
	m.Url = url
	return m
}

func (m MongoDatabase) Save(ctx context.Context, collectionName string, dto interface{}) error {
	_, err := m.Client.Database(m.Name).Collection(collectionName).InsertOne(ctx, dto)
	return err
}

func (m MongoDatabase) Update(ctx context.Context, collectionName string, filter, dto interface{}) error {
	_, err := m.Client.Database(m.Name).Collection(collectionName).UpdateOne(ctx, filter, dto)
	return err
}

func (m MongoDatabase) SaveMany(ctx context.Context, collectionName string, dtos []interface{}) error {
	_, err := m.Client.Database(m.Name).Collection(collectionName).InsertMany(ctx, dtos)
	return err
}

func (m MongoDatabase) GetById(ctx context.Context, collectionName string, filter interface{}, dto interface{}) error {
	err := m.Client.Database(m.Name).Collection(collectionName).FindOne(ctx, filter).Decode(dto)
	return err
}

func (m MongoDatabase) GetByFilter(ctx context.Context, collectionName string, filter interface{}, opt *options.FindOptions, dto interface{}) ([]byte, error) {
	results := make([]interface{}, 0)

	objectType := reflect.TypeOf(dto).Elem()

	cur, err := m.Client.Database(m.Name).Collection(collectionName).Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		result := reflect.New(objectType).Interface()
		err := cur.Decode(result)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return json.Marshal(results)
}

func (m MongoDatabase) Delete(ctx context.Context, collectionName string, filter interface{}) error {
	_, err := m.Client.Database(m.Name).Collection(collectionName).DeleteOne(ctx, filter)
	return err
}

func (m MongoDatabase) DeleteMany(ctx context.Context, collectionName string, filter interface{}) error {
	_, err := m.Client.Database(m.Name).Collection(collectionName).DeleteMany(ctx, filter)
	return err
}
