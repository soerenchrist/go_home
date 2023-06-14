package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/soerenchrist/mini_home/models"
)

type SqliteDevicesDatabase struct {
	db *sql.DB
}

func NewDevicesDatabase(db *sql.DB) (*SqliteDevicesDatabase, error) {
	database := &SqliteDevicesDatabase{db: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (db *SqliteDevicesDatabase) Close() error {
	return db.db.Close()
}

func (db *SqliteDevicesDatabase) AddDevice(device models.Device) error {
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

func (db *SqliteDevicesDatabase) GetDevice(id string) (models.Device, error) {
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

func (db *SqliteDevicesDatabase) DeleteDevice(id string) error {
	stmt, err := db.db.Prepare("delete from devices where id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return &models.NotFoundError{Message: "Device not found"}
	}

	return nil
}

func (db *SqliteDevicesDatabase) ListDevices() ([]models.Device, error) {
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

func (db *SqliteDevicesDatabase) GetSensor(deviceId, sensorId string) (models.Sensor, error) {
	stmt, err := db.db.Prepare("select id, name, data_type, device_id, sensor_type, is_active, unit, polling_interval from sensors where id = ? and device_id = ?")
	if err != nil {
		return models.Sensor{}, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(sensorId, deviceId)
	var sensor models.Sensor
	sensor, err = readSensor(row)
	if err != nil {
		return models.Sensor{}, err
	}

	return sensor, nil
}

func readSensor(row interface {
	Scan(dest ...interface{}) error
}) (models.Sensor, error) {
	var id string
	var name string
	var sensorType models.SensorType
	var isActive bool
	var dataType models.DataType
	var deviceId string
	var unit string
	var pollingInterval int
	err := row.Scan(&id, &name, &dataType, &deviceId, &sensorType, &isActive, &unit, &pollingInterval)

	if err != nil {
		return models.Sensor{}, err
	}
	return models.Sensor{ID: id, Name: name, DataType: dataType, DeviceID: deviceId, Type: sensorType, IsActive: isActive, Unit: unit, PollingInterval: pollingInterval}, nil
}

func (db *SqliteDevicesDatabase) DeleteSensor(deviceId, sensorId string) error {
	stmt, err := db.db.Prepare("delete from sensors where id = ? and device_id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(sensorId, deviceId)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return &models.NotFoundError{Message: "Sensor not found"}
	}

	return nil
}

func (db *SqliteDevicesDatabase) ListSensors(deviceId string) ([]models.Sensor, error) {
	rows, err := db.db.Query("select id, name, data_type, device_id, sensor_type, is_active, unit, polling_interval from sensors where device_id = ?", deviceId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]models.Sensor, 0)
	for rows.Next() {
		sensor, err := readSensor(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, sensor)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *SqliteDevicesDatabase) ListPollingSensors() ([]models.Sensor, error) {
	rows, err := db.db.Query("select id, name, data_type, device_id, sensor_type, is_active, unit, polling_interval from sensors where sensor_type = 'polling'")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]models.Sensor, 0)
	for rows.Next() {
		sensor, err := readSensor(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, sensor)
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

	stmt, err := tx.Prepare("insert into sensors(id, name, device_id, data_type, sensor_type, is_active, unit, polling_interval) values(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(sensor.ID, sensor.Name, sensor.DeviceID, sensor.DataType, sensor.Type, sensor.IsActive, sensor.Unit, sensor.PollingInterval)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (db *SqliteDevicesDatabase) AddSensorValue(data models.SensorValue) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into sensor_values(sensor_id, device_id, timestamp, value) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(data.SensorID, data.DeviceID, data.Timestamp, data.Value); err != nil {
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
	if _, err := db.db.Exec(createDevicesTableStmt); err != nil {
		return err
	}

	createSensorsTableStmt := `
	create table if not exists sensors (
		id text not null primary key,
		device_id text not null,
		name text,
		data_type text,
		unit text,
		sensor_type text,
		is_active integer,
		polling_interval integer,
		foreign key(device_id) references devices(id) on delete cascade
	);
	`

	if _, err := db.db.Exec(createSensorsTableStmt); err != nil {
		return err
	}

	createSensorDataTableStmt := `
	create table if not exists sensor_values (
		sensor_id text not null,
		device_id text not null,
		timestamp text not null,
		value text not null,
		foreign key(sensor_id) references sensors(id) on delete cascade,
		foreign key(device_id) references devices(id) on delete cascade
	);`

	if _, err := db.db.Exec(createSensorDataTableStmt); err != nil {
		return err
	}
	return nil
}

func (database *SqliteDevicesDatabase) SeedDatabase() {
	device1 := models.Device{ID: "1", Name: "My Device 1"}
	sensor1 := models.Sensor{ID: "S1", Name: "Temperature", DeviceID: "1", DataType: models.DataTypeFloat, Type: models.SensorTypeExternal, IsActive: true, Unit: "Celsius", PollingInterval: 0}
	sensor2 := models.Sensor{ID: "S2", Name: "Availability", DeviceID: "1", DataType: models.DataTypeBool, Type: models.SensorTypePolling, IsActive: true, Unit: "", PollingInterval: 10}

	device2 := models.Device{ID: "2", Name: "My Device 2"}
	sensor3 := models.Sensor{ID: "S3", Name: "Filling Level", DeviceID: "2", DataType: models.DataTypeInt, Type: models.SensorTypeExternal, IsActive: true, Unit: "%", PollingInterval: 0}

	database.AddDevice(device1)
	database.AddDevice(device2)

	database.AddSensor(sensor1)
	database.AddSensor(sensor2)
	database.AddSensor(sensor3)
}
