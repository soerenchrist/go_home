package tests

import (
	"encoding/json"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/rules"
)

func TestListRules_ShouldReturnRules(t *testing.T) {
	w := RecordGetCall(t, "/api/v1/rules")

	assert.Equal(t, w.Code, 200)

	var results []rules.Rule
	err := json.Unmarshal(w.Body.Bytes(), &results)
	if err != nil {
		t.Errorf("Error while unmarshalling rules: %s", err.Error())
		return
	}

	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Id, int64(1))
	assert.Equal(t, results[0].Name, "Turn on light when temperature is below 20")
	assert.Equal(t, results[0].When, rules.WhenExpression("when ${1.S1.current} < 20"))
	assert.Equal(t, results[0].Then, rules.ThenExpression("then ${1.C1} params {\"p_payload\": \"on\"}"))
}

func TestPostRule_ShouldReturn400_WhenJsonIsInvalid(t *testing.T) {
	w := RecordPostCall(t, "/api/v1/rules", "invalid")

	assert.Equal(t, w.Code, 400)
}

func TestPostRule_ShouldAddToDatabase(t *testing.T) {
	validator := func(database db.DevicesDatabase) {
		results, err := database.ListRules()
		if err != nil {
			t.Errorf("Error while listing rules: %s", err.Error())
			return
		}

		assert.Equal(t, len(results), 2)
		assert.Equal(t, results[1].Id, int64(2))
		assert.Equal(t, results[1].Name, "Test")
		assert.Equal(t, results[1].When, rules.WhenExpression("When"))
		assert.Equal(t, results[1].Then, rules.ThenExpression("Then"))
	}

	body := `
	{
		"name": "Test",
		"when": "When",
		"then": "Then"
	}
	`

	w := RecordPostCallWithDb(t, "/api/v1/rules", body, validator)

	assert.Equal(t, w.Code, 201)

	var rule rules.Rule
	err := json.Unmarshal(w.Body.Bytes(), &rule)
	if err != nil {
		t.Errorf("Error while unmarshalling rule: %s", err.Error())
		return
	}

	assert.Equal(t, rule.Id, int64(2))
	assert.Equal(t, rule.Name, "Test")
	assert.Equal(t, rule.When, rules.WhenExpression("When"))
	assert.Equal(t, rule.Then, rules.ThenExpression("Then"))
}
