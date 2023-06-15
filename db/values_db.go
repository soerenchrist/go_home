package db

import (
	"time"

	"github.com/soerenchrist/go_home/models"
)

func (db *SqliteDevicesDatabase) AddSensorValue(data *models.SensorValue) error {
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

func (db *SqliteDevicesDatabase) GetCurrentSensorValue(deviceId, sensorId string) (*models.SensorValue, error) {
	row := db.db.QueryRow("select sensor_id, device_id, timestamp, value from sensor_values where sensor_id = ? and device_id = ? order by timestamp desc limit 1", sensorId, deviceId)

	var sid string
	var did string
	var timestamp string
	var value string
	if err := row.Scan(&sid, &did, &timestamp, &value); err != nil {
		return nil, err
	}

	return &models.SensorValue{
		SensorID:  sid,
		DeviceID:  did,
		Timestamp: timestamp,
		Value:     value,
	}, nil
}

func (db *SqliteDevicesDatabase) GetSensorValuesSince(deviceId, sensorId string, timestamp time.Time) ([]models.SensorValue, error) {
	time_str := timestamp.Format(time.RFC3339)
	rows, err := db.db.Query("select sensor_id, device_id, timestamp, value from sensor_values where timestamp > ? and device_id = ? and sensor_id = ?", time_str, deviceId, sensorId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]models.SensorValue, 0)
	for rows.Next() {
		var sid string
		var did string
		var timestamp string
		var value string
		if err := rows.Scan(&sid, &did, &timestamp, &value); err != nil {
			return nil, err
		}
		results = append(results, models.SensorValue{
			SensorID:  sid,
			DeviceID:  did,
			Timestamp: timestamp,
			Value:     value,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (database *SqliteDevicesDatabase) createSensorValuesTable() error {
	createSensorDataTableStmt := `
	create table if not exists sensor_values (
		sensor_id text not null,
		device_id text not null,
		timestamp text not null,
		value text not null,
		foreign key(sensor_id) references sensors(id) on delete cascade,
		foreign key(device_id) references devices(id) on delete cascade
	);`

	if _, err := database.db.Exec(createSensorDataTableStmt); err != nil {
		return err
	}
	return nil
}
