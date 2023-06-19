package db

import (
	"github.com/soerenchrist/go_home/models"
)

func (db *SqliteDevicesDatabase) AddDevice(device *models.Device) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into devices(id, name) values(?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(device.ID, device.Name)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (db *SqliteDevicesDatabase) GetDevice(id string) (*models.Device, error) {
	stmt, err := db.db.Prepare("select id, name from devices where id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var deviceId string
	var name string
	err = stmt.QueryRow(id).Scan(&deviceId, &name)
	if err != nil {
		return nil, err
	}

	return &models.Device{ID: deviceId, Name: name}, nil
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
	rows, err := db.db.Query("select id, name from devices")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]models.Device, 0)
	for rows.Next() {
		var id string
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		results = append(results, models.Device{ID: id, Name: name})
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
		name text
	);
	`
	if _, err := database.db.Exec(createDevicesTableStmt); err != nil {
		return err
	}

	return nil
}
