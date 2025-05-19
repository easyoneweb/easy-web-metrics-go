package metrics

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ikirja/easy-web-metrics-go/internal/database"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type utm struct {
	UtmSource   string `json:"utmSource"`
	UtmMedium   string `json:"utmMedium"`
	UtmCampaign string `json:"utmCampaign"`
}

type user struct {
	UserID     string `json:"userID"`
	Login      string `json:"login"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	LastName   string `json:"lastName"`
	Phone      string `json:"phone"`
}

type visitor struct {
	Visitor   string `json:"visitor"`
	Url       string `json:"url"`
	Utm       utm    `json:"utm"`
	Referrer  string `json:"referrer"`
	IP        string `json:"ip"`
	UserAgent string `json:"userAgent"`
	UserData  user   `json:"userData"`
}

type processedVisitor struct {
	CreatedAt time.Time `json:"createdAt"`
	Visitor   string    `json:"visitor"`
	Status    string    `json:"status"`
}

type visitorStatus struct {
	New     string
	Updated string
}

var vStatus = visitorStatus{
	New:     "new",
	Updated: "updated",
}

func ProcessVisitor(r *http.Request) (processedVisitor, error) {
	var err error

	v := visitor{}
	pVisitor := processedVisitor{
		Visitor: "",
		Status:  "",
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&v)
	if err != nil {
		return pVisitor, errors.New("check sent data")
	}

	filter := getBsonFilter(v)

	urlDB := database.UrlDB{
		Url:      v.Url,
		Utm:      database.UtmDB(v.Utm),
		Referrer: v.Referrer,
	}
	visitorDB := database.VisitorDB{
		Visitor:   v.Visitor,
		Urls:      []database.UrlDB{urlDB},
		IP:        v.IP,
		UserAgent: v.UserAgent,
		UserData:  database.UserDB(v.UserData),
	}

	updatedVisitor, _ := database.VisitorUpdate(visitorDB, filter)
	if updatedVisitor.Visitor != "" {
		pVisitor.CreatedAt = updatedVisitor.CreatedAt
		pVisitor.Visitor = updatedVisitor.Visitor
		pVisitor.Status = vStatus.Updated
		return pVisitor, nil
	}

	createdVisitor, err := database.VisitorCreate(visitorDB)
	if err != nil {
		return pVisitor, err
	}

	pVisitor.CreatedAt = createdVisitor.CreatedAt
	pVisitor.Visitor = createdVisitor.Visitor
	pVisitor.Status = vStatus.New
	return pVisitor, nil
}

func GetVisitors(limit int64, skip int64) ([]database.VisitorDB, error) {
	return database.GetVisitors(limit, skip)
}

func getBsonFilter(v visitor) bson.D {
	var filter bson.D

	if v.UserData.UserID != "" {
		filter = bson.D{{"userData.userID", v.UserData.UserID}}
		return filter
	}
	if v.Visitor != "" {
		filter = bson.D{{"visitor", v.Visitor}}
		return filter
	}
	if v.Visitor == "" && v.UserData.UserID == "" && v.UserAgent != "" && v.IP != "" {
		filter = bson.D{{"userAgent", v.UserAgent}, {"ip", v.IP}}
		return filter
	}

	return filter
}
