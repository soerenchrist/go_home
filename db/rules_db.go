package db

import "github.com/soerenchrist/go_home/rules"

func (database *SqliteDevicesDatabase) AddRule(rule *rules.Rule) error {
	result := database.db.Create(rule)
	return result.Error
}

func (database *SqliteDevicesDatabase) ListRules() ([]rules.Rule, error) {
	rules := make([]rules.Rule, 0)
	result := database.db.Find(&rules)
	return rules, result.Error
}
