package messages

type errors struct {
	Test test
}

type test struct {
	DB      db
	Visitor visitor
}

type db struct {
	Connect       string
	VisitorCreate string
	VisitorID     string
	VisitorUpdate string
	VisitorDelete string
}

type visitor struct {
	ProcessVisitor string
	VisitorID      string
}

var Errors = errors{
	Test: test{
		DB: db{
			Connect:       "couldn't connect to db",
			VisitorCreate: "couldn't create visitor",
			VisitorID:     "visitor id is not uuid",
			VisitorUpdate: "couldn't update visitor",
			VisitorDelete: "couldn't delete visitor",
		},
		Visitor: visitor{
			ProcessVisitor: "couldn't process visitor",
			VisitorID:      "visitor id is not uuid",
		},
	},
}
