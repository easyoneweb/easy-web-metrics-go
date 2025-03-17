package database

type UtmDB struct {
	UtmSource   string `bson:"utmSource"`
	UtmMedium   string `bson:"utmMedium"`
	UtmCampaign string `bson:"utmCampaign"`
}

type UrlDB struct {
	Url      string `bson:"url"`
	Utm      UtmDB  `bson:"utm"`
	Referrer string `bson:"referrer"`
}

type UserDB struct {
	UserID   string `bson:"userID"`
	Login      string `bson:"login"`
	Email      string `bson:"email"`
	FirstName  string `bson:"firstName"`
	SecondName string `bson:"secondName"`
	LastName   string `bson:"lastName"`
	Phone      string `bson:"phone"`
}

type VisitorDB struct {
	Visitor   string  `bson:"visitor"`
	Urls      []UrlDB `bson:"url"`
	IP        string  `bson:"ip"`
	UserAgent string  `bson:"userAgent"`
	UserData  UserDB  `bson:"userData"`
}
