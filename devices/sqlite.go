package devices

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDatabase struct {
	db   *sql.DB
	path string
}

func NewSqliteDatabase(path string) (*SqliteDatabase, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	database := &SqliteDatabase{path: path, db: db}
	err = database.createTables()
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (db *SqliteDatabase) Close() error {
	return db.db.Close()
}

func (db *SqliteDatabase) Add(device Device) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into devices(id, name, last_reached) values(?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(device.ID, device.Name, device.LastReached)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (db *SqliteDatabase) List() ([]Device, error) {
	rows, err := db.db.Query("select id, name, last_reached from devices")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]Device, 0)
	for rows.Next() {
		var id string
		var name string
		var lastReached string
		err = rows.Scan(&id, &name, &lastReached)
		if err != nil {
			return nil, err
		}
		results = append(results, Device{ID: id, Name: name, LastReached: lastReached})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *SqliteDatabase) createTables() error {
	sqlStmt := `
	create table if not exists devices (
		id text not null primary key,
		name text,
		last_reached text
	);
	`
	_, err := db.db.Exec(sqlStmt)
	return err
}
