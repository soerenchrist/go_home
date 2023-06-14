package db

import "github.com/soerenchrist/mini_home/models"

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

func (database *SqliteDevicesDatabase) createDeviceTable() error {
	createDevicesTableStmt := `
	create table if not exists devices (
		id text not null primary key,
		name text,
		last_reached text
	);
	`
	if _, err := database.db.Exec(createDevicesTableStmt); err != nil {
		return err
	}

	return nil
}
