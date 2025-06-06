package database

import (
	"context"
	"errors"
	"log"
	"time"

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
	createVisitorsIndexes(collectionVisitors)
	return nil
}

func getCollection(collection string) *mongo.Collection {
	return db.Collection(collection)
}

func createVisitorsIndexes(collectionVisitors *mongo.Collection) {
	indexModels := []mongo.IndexModel{
		{Keys: bson.D{{Key: "updatedAt", Value: -1}}},
		{Keys: bson.D{{Key: "userData.userID", Value: -1}}},
		{Keys: bson.D{{Key: "visitor", Value: 1}}},
		{Keys: bson.D{{Key: "ip", Value: 1}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := collectionVisitors.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Fatal(err)
	}
}

func VisitorCreate(v VisitorDB) (VisitorDB, error) {
	result := VisitorDB{}
	v.Visitor = uuid.NewString()
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()

	_, err := collectionVisitors.InsertOne(context.TODO(), v)
	if err != nil {
		return VisitorDB{}, errors.New("couldn't create new visitor")
	}

	collectionVisitors.FindOne(context.TODO(), bson.D{{"visitor", v.Visitor}}).Decode(&result)
	return result, nil
}

func VisitorUpdate(v VisitorDB, filter bson.D) (VisitorDB, error) {
	if filter == nil {
		return VisitorDB{}, errors.New("filter is empty")
	}

	result := VisitorDB{}

	err := collectionVisitors.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return VisitorDB{}, errors.New("couldn't find visitor")
		}
		log.Fatal(err)
	}

	if v.Visitor != result.Visitor {
		go deleteVisitorWithoutUser(v)
	}

	updateVisitorData := VisitorDB{
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  time.Now(),
		VisitDates: updateVisitorDates(result.VisitDates),
		Visitor:    result.Visitor,
		Urls:       updateVisitorUrls(result.Urls, v.Urls[0]),
		IP:         v.IP,
		UserAgent:  v.UserAgent,
		UserData:   updateVisitorUserData(result.UserData, v.UserData),
	}

	update := bson.D{{"$set", updateVisitorData}}

	_, err = collectionVisitors.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return VisitorDB{}, errors.New("couldn't update visitor")
	}

	return updateVisitorData, nil
}

func GetVisitors(limit int64, skip int64) ([]VisitorDB, error) {
	filter := bson.M{"updatedAt": bson.M{"$gte": time.Now().Add(-time.Hour * 720), "$lt": time.Now()}}

	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := collectionVisitors.Find(context.TODO(), filter, opts)
	if err != nil {
		return []VisitorDB{}, errors.New("couldn't execute collection find")
	}

	var results []VisitorDB
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []VisitorDB{}, errors.New("couldn't do get visitors cursor all")
	}

	return results, nil
}

func updateVisitorDates(visitDates []time.Time) []time.Time {
	if len(visitDates) == 0 {
		return []time.Time{time.Now()}
	}

	format := "2006-01-02"
	today := time.Now().Format(format)
	lastVisitDate := visitDates[len(visitDates)-1].Format(format)
	if today == lastVisitDate {
		return visitDates
	}

	const maxVisitDates = 30

	updatedVisitDates := make([]time.Time, 0)

	for i := len(visitDates) - 1; i >= 0; i-- {
		if len(updatedVisitDates) >= maxVisitDates-1 {
			break
		}
		updatedVisitDates = append([]time.Time{visitDates[i]}, updatedVisitDates...)
	}

	updatedVisitDates = append(updatedVisitDates, time.Now())
	return updatedVisitDates
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

func updateVisitorUserData(userDB UserDB, userToUpdate UserDB) UserDB {
	if userDB.UserID == "" && userToUpdate.UserID != "" {
		userDB = userToUpdate
	}

	return userDB
}
