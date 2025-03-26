package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ikirja/easy-web-metrics-go/internal/database"
	"github.com/ikirja/easy-web-metrics-go/internal/messages"

	"github.com/google/uuid"
)

func TestProcessVisitor(t *testing.T) {
	var dbName = "easywebmetricstest"

	err := database.Connect(dbName)
	if err != nil {
		t.Errorf("%v: %v", messages.Errors.Test.DB.Connect, err)
	}

	t.Run("Visitor create", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/api/v1/metrics/visitor", createBodyReader())
		processVisitor(t, r)
	})
	t.Run("Visitor update", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/api/v1/metrics/visitor", createBodyReader())
		processVisitor(t, r)
	})
	t.Run("Visitor get", func(t *testing.T) {
		_, err := GetVisitors()
		if err != nil {
			t.Errorf("%v: %v", messages.Errors.Test.Visitor.GetVisitors, err)
		}
	})
}

func createBodyReader() *strings.Reader {
	return strings.NewReader(`
		{
			"visitor": "",
			"url": "/some/test/url",
			"utm": {
				"utmSource": "test",
				"utmMedium": "test",
				"utmCampaign": "test"
			},
			"userData": {
					"userID": "8"
			}
		}
	`)
}

func processVisitor(t *testing.T, r *http.Request) {
	processedVisitor, err := ProcessVisitor(r)
	if err != nil {
		t.Errorf("%v: %v", messages.Errors.Test.Visitor.ProcessVisitor, err)
	}

	_, err = uuid.Parse(processedVisitor.Visitor)
	if err != nil {
		t.Errorf("%v: %v", messages.Errors.Test.DB.VisitorID, err)
	}
}
