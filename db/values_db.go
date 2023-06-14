package db

import "github.com/soerenchrist/mini_home/models"

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
