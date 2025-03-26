package database

import "time"

type UtmDB struct {
	UtmSource   string `json:"utmSource" bson:"utmSource"`
	UtmMedium   string `json:"utmMedium" bson:"utmMedium"`
	UtmCampaign string `json:"utmCampaign" bson:"utmCampaign"`
}

type UrlDB struct {
	Url      string `json:"url" bson:"url"`
	Utm      UtmDB  `json:"utm" bson:"utm"`
	Referrer string `json:"referrer" bson:"referrer"`
}

type UserDB struct {
	UserID     string `json:"userID" bson:"userID"`
	Login      string `json:"login" bson:"login"`
	Email      string `json:"email" bson:"email"`
	FirstName  string `json:"firstName" bson:"firstName"`
	SecondName string `json:"secondName" bson:"secondName"`
	LastName   string `json:"lastName" bson:"lastName"`
	Phone      string `json:"phone" bson:"phone"`
}

type VisitorDB struct {
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt" bson:"updatedAt"`
	VisitDates []time.Time `json:"visitDates" bson:"visitDates"`
	Visitor    string      `json:"visitor" bson:"visitor"`
	Urls       []UrlDB     `json:"url" bson:"url"`
	IP         string      `json:"ip" bson:"ip"`
	UserAgent  string      `json:"userAgent" bson:"userAgent"`
	UserData   UserDB      `json:"userData" bson:"userData"`
}
