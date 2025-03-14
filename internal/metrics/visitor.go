package metrics

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ikirja/easy-web-metrics-go/internal/database"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type utm struct {
	UtmSource   string `json:"utmSource"`
	UtmMedium   string `json:"utmMedium"`
	UtmCampaign string `json:"utmCampaign"`
}

type user struct {
	BitrixID   string `json:"bitrixID"`
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
	Visitor string `json:"visitor"`
	Status  string `json:"status"`
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
	var filter bson.D

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

	v.IP = readIpAddress(r)
	v.UserAgent = r.Header.Get("user-agent")

	if validateVisitorData(v, "visitor") {
		filter = bson.D{{"visitor", v.Visitor}}
	}

	if validateVisitorData(v, "bitrixID") {
		filter = bson.D{{"userData.bitrixID", v.UserData.BitrixID}}
	}

	if validateVisitorData(v, "userAgent") {
		filter = bson.D{{"userAgent", v.UserAgent}, {"ip", v.IP}}
	}

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
		pVisitor.Visitor = updatedVisitor.Visitor
		pVisitor.Status = vStatus.Updated
		return pVisitor, nil
	}

	createdVisitor, err := database.VisitorCreate(visitorDB)
	if err != nil {
		return pVisitor, err
	}

	pVisitor.Visitor = createdVisitor.Visitor
	pVisitor.Status = vStatus.New
	return pVisitor, nil
}

func readIpAddress(r *http.Request) string {
	ip := r.Header.Get("x-real-ip")
	if ip == "" {
		ip = r.Header.Get("x-forwarded-for")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func validateVisitorData(v visitor, visitorDataType string) bool {
	if visitorDataType == "visitor" {
		if v.Visitor != "" {
			return true
		}
	}
	if visitorDataType == "bitrixID" {
		if v.Visitor == "" && v.UserData.BitrixID != "" {
			return true
		}
	}
	if visitorDataType == "userAgent" {
		if v.Visitor == "" && v.UserData.BitrixID == "" && v.UserAgent != "" && v.IP != "" {
			return true
		}
	}
	return false
}
