package database

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DBConfig struct {
	DB                 *mongo.Database
	CollectionVisitors *mongo.Collection
}

const collectionVisitorsName = "visitors"

var db *mongo.Database
var collectionVisitors *mongo.Collection

func Connect(dbName string) error {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	db = client.Database(dbName)
	collectionVisitors = getCollection(collectionVisitorsName)
	return nil
}

func getCollection(collection string) *mongo.Collection {
	return db.Collection(collection)
}

func VisitorCreate(v VisitorDB) (VisitorDB, error) {
	result := VisitorDB{}
	v.Visitor = uuid.NewString()

	_, err := collectionVisitors.InsertOne(context.TODO(), v)
	if err != nil {
		return VisitorDB{}, errors.New("couldn't create new visitor")
	}

	collectionVisitors.FindOne(context.TODO(), bson.D{{"visitor", v.Visitor}}).Decode(&result)
	return result, nil
}

func VisitorUpdate(v VisitorDB, filter bson.D) (VisitorDB, error) {
	result := VisitorDB{}

	err := collectionVisitors.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return VisitorDB{}, errors.New("couldn't find visitor")
		}
		log.Fatal(err)
	}

	// UPDATE ONLY WHAT IS NEEDED TO BE UPDATED!
	updateVisitorData := VisitorDB{
		Visitor:   result.Visitor,
		Urls:      updateVisitorUrls(result.Urls, v.Urls[0]),
		IP:        v.IP,
		UserAgent: v.UserAgent,
		UserData:  v.UserData,
	}

	update := bson.D{{"$set", updateVisitorData}}

	_, err = collectionVisitors.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return VisitorDB{}, errors.New("couldn't update visitor")
	}

	return updateVisitorData, nil
}

func updateVisitorUrls(urls []UrlDB, newUrl UrlDB) []UrlDB {
	const maxUrls = 10

	updatedUrls := make([]UrlDB, 0)

	for i := len(urls) - 1; i >= 0; i-- {
		if len(updatedUrls) >= maxUrls-1 {
			break
		}
		updatedUrls = append([]UrlDB{urls[i]}, updatedUrls...)
	}

	updatedUrls = append(updatedUrls, newUrl)
	return updatedUrls
}
