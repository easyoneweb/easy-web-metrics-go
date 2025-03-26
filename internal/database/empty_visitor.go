package database

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func deleteVisitorWithoutUser(v VisitorDB) error {
	if v.Visitor == "" {
		return nil
	}

	result := VisitorDB{}

	err := collectionVisitors.FindOne(context.TODO(), bson.D{{"visitor", v.Visitor}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		log.Fatal(err)
	}

	if isEmptyUser(result) {
		_, err := collectionVisitors.DeleteOne(context.TODO(), bson.D{{"visitor", result.Visitor}})
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func isEmptyUser(result VisitorDB) bool {
	return result.UserData.UserID == ""
}
