package db

import "github.com/soerenchrist/go_home/models"

func (db *SqliteDevicesDatabase) GetSensor(deviceId, sensorId string) (*models.Sensor, error) {
	stmt, err := db.db.Prepare("select id, name, data_type, device_id, sensor_type, is_active, unit, polling_interval, polling_endpoint, polling_strategy from sensors where id = ? and device_id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(sensorId, deviceId)
	var sensor *models.Sensor
	sensor, err = readSensor(row)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

func (db *SqliteDevicesDatabase) ListSensors(deviceId string) ([]models.Sensor, error) {
	rows, err := db.db.Query("select id, name, data_type, device_id, sensor_type, is_active, unit, polling_interval, polling_endpoint, polling_strategy from sensors where device_id = ?", deviceId)
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
		results = append(results, *sensor)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *SqliteDevicesDatabase) ListPollingSensors() ([]models.Sensor, error) {
	rows, err := db.db.Query("select id, name, data_type, device_id, sensor_type, is_active, unit, polling_interval, polling_endpoint, polling_strategy from sensors where sensor_type = 'polling'")
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
		results = append(results, *sensor)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *SqliteDevicesDatabase) AddSensor(sensor *models.Sensor) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into sensors(id, name, device_id, data_type, sensor_type, is_active, unit, polling_interval, polling_endpoint, polling_strategy) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(sensor.ID, sensor.Name, sensor.DeviceID, sensor.DataType, sensor.Type, sensor.IsActive, sensor.Unit, sensor.PollingInterval, sensor.PollingEndpoint, sensor.PollingStrategy)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
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

func readSensor(row interface {
	Scan(dest ...interface{}) error
}) (*models.Sensor, error) {
	var id string
	var name string
	var sensorType models.SensorType
	var isActive bool
	var dataType models.DataType
	var deviceId string
	var unit string
	var pollingInterval int
	var pollingEndpoint string
	var pollingStrategy models.PollingStrategy
	err := row.Scan(&id, &name, &dataType, &deviceId, &sensorType, &isActive, &unit, &pollingInterval, &pollingEndpoint, &pollingStrategy)

	if err != nil {
		return nil, err
	}
	return &models.Sensor{
		ID:              id,
		Name:            name,
		DataType:        dataType,
		DeviceID:        deviceId,
		Type:            sensorType,
		IsActive:        isActive,
		Unit:            unit,
		PollingInterval: pollingInterval,
		PollingEndpoint: pollingEndpoint,
		PollingStrategy: pollingStrategy}, nil
}

func (database *SqliteDevicesDatabase) createSensorsTable() error {

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
		polling_endpoint text,
		polling_strategy text,
		foreign key(device_id) references devices(id) on delete cascade
	);
	`

	if _, err := database.db.Exec(createSensorsTableStmt); err != nil {
		return err
	}
	return nil
}
