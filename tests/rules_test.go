package tests

import (
	"encoding/json"
	"fmt"
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

func TestPostRule_ShouldReturn400WhenRuleIsInvalid(t *testing.T) {
	invalidRules := []rules.Rule{
		{
			Name: "",
			When: rules.WhenExpression("when ${1.S1.current} < 20"),
			Then: rules.ThenExpression("then ${1.C1}"),
		},
		{
			Name: "Test",
			When: rules.WhenExpression(""),
			Then: rules.ThenExpression("then ${1.C1}"),
		},
		{
			Name: "Test",
			When: rules.WhenExpression("when ${1.S1.current} < 20"),
			Then: rules.ThenExpression(""),
		},
	}

	expectedMessages := []string{
		"Name is required",
		"invalid rule: When Expression is empty",
		"invalid rule: Then Expression is empty",
	}

	for i, rule := range invalidRules {
		body, err := json.Marshal(rule)
		if err != nil {
			t.Errorf("Error while marshalling rule: %s", err.Error())
			return
		}

		w := RecordPostCall(t, "/api/v1/rules", string(body))

		assert.Equal(t, w.Code, 400)

		assertErrorMessageEquals(t, w.Body.Bytes(), expectedMessages[i])
	}
}

func TestPostRule_ShouldAddToDatabase(t *testing.T) {
	when := "when ${1.S1.current} < 20"
	validator := func(database db.DevicesDatabase) {
		results, err := database.ListRules()
		if err != nil {
			t.Errorf("Error while listing rules: %s", err.Error())
			return
		}

		assert.Equal(t, len(results), 2, "Is not added to database")
		assert.Equal(t, results[1].Id, int64(2))
		assert.Equal(t, results[1].Name, "Test")
		assert.Equal(t, results[1].When, rules.WhenExpression(when))
		assert.Equal(t, results[1].Then, rules.ThenExpression("Then ${1.C1}"))
	}

	body := fmt.Sprintf(`
	{
		"name": "Test",
		"when": "%s",
		"then": "Then ${1.C1}"
	}
	`, when)

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
	assert.Equal(t, rule.When, rules.WhenExpression(when))
	assert.Equal(t, rule.Then, rules.ThenExpression("Then ${1.C1}"))
}
