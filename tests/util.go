package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/soerenchrist/go_home/internal/db"
	"github.com/soerenchrist/go_home/internal/server"
	"github.com/soerenchrist/go_home/internal/value"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateTestDatabase(filename string) db.Database {
	url := fmt.Sprintf("file:%s?mode=memory", filename)

	sqlite := sqlite.Open(url)

	gormDb, err := gorm.Open(sqlite, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database, err := db.NewDevicesDatabase(gormDb)
	if err != nil {
		panic(err)
	}

	database.SeedDatabase()

	return database
}

type DbValidator func(database db.Database)

func recordCall(t *testing.T, url string, method string, body io.Reader, dbValidator DbValidator) *httptest.ResponseRecorder {
	gin.DefaultWriter = io.Discard
	w := httptest.NewRecorder()
	filename := t.Name()
	database := CreateTestDatabase(filename)
	if dbValidator != nil {
		defer dbValidator(database)
	}
	outputBindings := value.NewOutputBindings()
	router := server.NewRouter(database, outputBindings)

	req := httptest.NewRequest(method, url, body)

	router.ServeHTTP(w, req)
	return w
}

func RecordGetCall(t *testing.T, url string) *httptest.ResponseRecorder {
	return recordCall(t, url, "GET", nil, nil)
}
func RecordDeleteCallWithDb(t *testing.T, url string, dbValidator DbValidator) *httptest.ResponseRecorder {
	return recordCall(t, url, "DELETE", nil, dbValidator)
}

func RecordDeleteCall(t *testing.T, url string) *httptest.ResponseRecorder {
	return recordCall(t, url, "DELETE", nil, nil)
}

func RecordPostCall(t *testing.T, url string, body string) *httptest.ResponseRecorder {
	reader := strings.NewReader(body)

	return recordCall(t, url, "POST", reader, nil)
}

func RecordPostCallWithDb(t *testing.T, url string, body string, dbValidator DbValidator) *httptest.ResponseRecorder {
	reader := strings.NewReader(body)

	return recordCall(t, url, "POST", reader, dbValidator)
}

func IsValidUuid(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

type errorResponse struct {
	Error string `json:"error"`
}

func assertErrorMessageEquals(t *testing.T, body []byte, expected string) {
	var data errorResponse
	err := json.Unmarshal(body, &data)
	if err != nil {
		t.Fatalf("Could not parse error message: %s", err)
	}

	if data.Error != expected {
		t.Fatalf("Expected error message '%s', got '%s'", expected, data.Error)
	}
}
