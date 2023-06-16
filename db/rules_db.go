package db

import "github.com/soerenchrist/go_home/rules"

func (database *SqliteDevicesDatabase) AddRule(rule *rules.Rule) error {
	stmt, err := database.db.Prepare("insert into rules (name, when_exp, then_exp) values (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(rule.Name, rule.When, rule.Then)
	if err != nil {
		return err
	}

	rule.Id, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (database *SqliteDevicesDatabase) ListRules() ([]rules.Rule, error) {
	rows, err := database.db.Query("select id, name, when_exp, then_exp from rules")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]rules.Rule, 0)
	for rows.Next() {
		var id int64
		var name string
		var when rules.WhenExpression
		var then rules.ThenExpression
		err := rows.Scan(&id, &name, &when, &then)
		if err != nil {
			return nil, err
		}
		rule := rules.Rule{
			Id:   id,
			Name: name,
			When: when,
			Then: then,
		}
		results = append(results, rule)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (database *SqliteDevicesDatabase) createRulesTable() error {
	createRulesTableStmt := `
        create table if not exists rules (
            id integer primary key autoincrement,
            name text not null,
            when_exp text not null,
            then_exp text not null
        );
    `
	if _, err := database.db.Exec(createRulesTableStmt); err != nil {
		return err
	}

	return nil
}
