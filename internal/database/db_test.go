package database

import (
	"testing"

	"github.com/ikirja/easy-web-metrics-go/internal/messages"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var testUrl = UrlDB{
	Url: "/some/test/url",
	Utm: UtmDB{
		UtmSource:   "test",
		UtmMedium:   "test",
		UtmCampaign: "test",
	},
	Referrer: "https://domain.com",
}

var testUser = UserDB{
	UserID:   "",
	Login:      "",
	Email:      "",
	FirstName:  "",
	SecondName: "",
	LastName:   "",
	Phone:      "",
}

var testVisitor = VisitorDB{
	Visitor:   "",
	Urls:      []UrlDB{testUrl},
	IP:        "127.0.0.1",
	UserAgent: "test/user-agent",
	UserData:  testUser,
}

func TestConnect(t *testing.T) {
	var dbName = "easywebmetricstest"

	err := Connect(dbName)
	if err != nil {
		t.Errorf("%v: %v", messages.Errors.Test.DB.Connect, err)
	}
}

func TestVisitor(t *testing.T) {
	var err error
	visitor := VisitorDB{}

	t.Run("Visitor Create", func(t *testing.T) {
		visitor, err = VisitorCreate(testVisitor)
		if err != nil {
			t.Errorf("%v: %v", messages.Errors.Test.DB.VisitorCreate, err)
		}

		_, err = uuid.Parse(visitor.Visitor)
		if err != nil {
			t.Errorf("%v: %v", messages.Errors.Test.DB.VisitorID, err)
		}
	})
	t.Run("Visitor Update", func(t *testing.T) {
		testVisitor.UserData.UserID = "8"
		filter := bson.D{{"visitor", visitor.Visitor}}

		_, err := VisitorUpdate(testVisitor, filter)
		if err != nil {
			t.Errorf("%v: %v", messages.Errors.Test.DB.VisitorUpdate, err)
		}
	})
	t.Run("Visitor Delete - Empty UserID", func(t *testing.T) {
		testVisitor.UserData.UserID = ""
		
		visitor, err = VisitorCreate(testVisitor)
		if err != nil {
			t.Errorf("%v: %v", messages.Errors.Test.DB.VisitorCreate, err)
		}

		_, err = uuid.Parse(visitor.Visitor)
		if err != nil {
			t.Errorf("%v: %v", messages.Errors.Test.DB.VisitorID, err)
		}

		testVisitor.Visitor = visitor.Visitor
		
		err := deleteVisitorWithoutUser(testVisitor)
		if err != nil {
				t.Errorf("%v: %v", messages.Errors.Test.DB.VisitorDelete, err)
		}
	})
	t.Run("Visitor Delete - Valid User Data", func(t *testing.T) {
		testVisitor.UserData.UserID = "8"
		err := deleteVisitorWithoutUser(testVisitor)
		if err != nil {
				t.Errorf("%v: %v", messages.Errors.Test.DB.VisitorDelete, err)
		}
	})
}
