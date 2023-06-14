package db

import "github.com/soerenchrist/mini_home/models"

func (db *SqliteDevicesDatabase) GetCommand(deviceId, commandId string) (models.Command, error) {
	row := db.db.QueryRow("select id, name, payload_template, endpoint, method from commands where id = ? and device_id = ?", commandId, deviceId)

	command, err := readCommand(row)
	if err != nil {
		return models.Command{}, err
	}

	return command, nil
}

func (db *SqliteDevicesDatabase) ListCommands(deviceId string) ([]models.Command, error) {
	rows, err := db.db.Query("select id, name, payload_template, endpoint, method from commands where device_id = ?", deviceId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]models.Command, 0)
	for rows.Next() {
		command, err := readCommand(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, command)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (db *SqliteDevicesDatabase) AddCommand(command models.Command) error {
	stmt, err := db.db.Prepare("insert into commands(id, device_id, name, payload_template, endpoint, method) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(command.ID, command.DeviceID, command.Name, command.PayloadTemplate, command.Endpoint, command.Method)
	if err != nil {
		return err
	}

	return nil
}

func (db *SqliteDevicesDatabase) DeleteCommand(deviceId, commandId string) error {
	stmt, err := db.db.Prepare("delete from commands where id = ? and device_id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(commandId, deviceId)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return &models.NotFoundError{Message: "Command not found"}
	}

	return nil
}

func readCommand(row interface {
	Scan(dest ...interface{}) error
}) (models.Command, error) {
	var id, name, payloadTemplate, endpoint, method string
	err := row.Scan(&id, &name, &payloadTemplate, &endpoint, &method)

	if err != nil {
		return models.Command{}, err
	}

	return models.Command{
		ID:              id,
		Name:            name,
		PayloadTemplate: payloadTemplate,
		Endpoint:        endpoint,
		Method:          method,
	}, nil
}

func (database *SqliteDevicesDatabase) createCommandsTable() error {
	createCommandsTableStmt := `
		create table if not exists commands (
			id text not null primary key,
			device_id text not null,
			name text not null,
			payload_template text not null,
			endpoint text not null,
			method text not null,
			foreign key(device_id) references devices(id) on delete cascade);`

	if _, err := database.db.Exec(createCommandsTableStmt); err != nil {
		return err
	}

	return nil
}
