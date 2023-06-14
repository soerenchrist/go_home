package tests

import (
	"database/sql"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/soerenchrist/mini_home/db"
	"github.com/soerenchrist/mini_home/server"
)

func CreateTestDatabase(filename string) db.DevicesDatabase {
	sqlite, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}
	if _, err := sqlite.Exec("PRAGMA foreign_keys = ON"); err != nil {
		panic(err)
	}

	database, err := db.NewDevicesDatabase(sqlite)
	if err != nil {
		panic(err)
	}

	database.SeedDatabase()

	return database
}

func CloseTestDatabase(database db.DevicesDatabase, filename string) {
	database.Close()
	os.Remove(filename)
}

type DbValidator func(database db.DevicesDatabase)

func recordCall(t *testing.T, url string, method string, body io.Reader, dbValidator DbValidator) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	filename := t.Name() + ".db"
	database := CreateTestDatabase(filename)
	defer CloseTestDatabase(database, filename)
	if dbValidator != nil {
		defer dbValidator(database)
	}
	router := server.NewRouter(database)

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
