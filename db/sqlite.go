package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/soerenchrist/mini_home/models"
)

type SqliteDevicesDatabase struct {
	db   *sql.DB
	path string
}

func NewDevicesDatabase(path string) (*SqliteDevicesDatabase, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	database := &SqliteDevicesDatabase{path: path, db: db}
	err = database.createTables()
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (db *SqliteDevicesDatabase) Close() error {
	return db.db.Close()
}

func (db *SqliteDevicesDatabase) Add(device models.Device) error {
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

func (db *SqliteDevicesDatabase) Get(id string) (models.Device, error) {
	stmt, err := db.db.Prepare("select id, name, last_reached from devices where id = ?")
	if err != nil {
		return models.Device{}, err
	}

	defer stmt.Close()

	var deviceId string
	var name string
	var lastReached string
	err = stmt.QueryRow(id).Scan(&deviceId, &name, &lastReached)
	if err != nil {
		return models.Device{}, err
	}

	return models.Device{ID: deviceId, Name: name, LastReached: lastReached}, nil
}

func (db *SqliteDevicesDatabase) List() ([]models.Device, error) {
	rows, err := db.db.Query("select id, name, last_reached from devices")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]models.Device, 0)
	for rows.Next() {
		var id string
		var name string
		var lastReached string
		err = rows.Scan(&id, &name, &lastReached)
		if err != nil {
			return nil, err
		}
		results = append(results, models.Device{ID: id, Name: name, LastReached: lastReached})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *SqliteDevicesDatabase) ListSensors(deviceId string) ([]models.Sensor, error) {
	rows, err := db.db.Query("select id, name, data_type, device_id, unit from sensors where device_id = ?", deviceId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]models.Sensor, 0)
	for rows.Next() {
		var id string
		var name string
		var dataType models.DataType
		var deviceId string
		var unit string
		err = rows.Scan(&id, &name, &dataType, &deviceId, &unit)
		if err != nil {
			return nil, err
		}
		results = append(results, models.Sensor{ID: id, Name: name, DataType: dataType, DeviceID: deviceId})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *SqliteDevicesDatabase) AddSensor(sensor models.Sensor) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into sensors(id, name, device_id, data_type) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(sensor.ID, sensor.Name, sensor.DeviceID, sensor.DataType)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (db *SqliteDevicesDatabase) createTables() error {
	createDevicesTableStmt := `
	create table if not exists devices (
		id text not null primary key,
		name text,
		last_reached text
	);
	`
	_, err := db.db.Exec(createDevicesTableStmt)

	if err != nil {
		return err
	}

	createSensorsTableStmt := `
	create table if not exists sensors (
		id text not null primary key,
		device_id text not null,
		name text,
		data_type text,
		unit text,
		foreign key(device_id) references devices(id)
	);
	`

	_, err = db.db.Exec(createSensorsTableStmt)
	return err
}
