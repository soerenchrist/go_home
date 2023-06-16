package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/db"
	"github.com/soerenchrist/go_home/models"
	"github.com/soerenchrist/go_home/rules"
)

type RulesController struct {
	database db.DevicesDatabase
}

func NewRulesController(database db.DevicesDatabase) *RulesController {
	return &RulesController{database: database}
}

func (controller *RulesController) ListRules(context *gin.Context) {
	rules, err := controller.database.ListRules()
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, rules)
}

func (controller *RulesController) PostRule(context *gin.Context) {
	var request models.CreateRuleRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rule := rules.Rule{
		Name: request.Name,
		When: rules.WhenExpression(request.When),
		Then: rules.ThenExpression(request.Then),
	}

	if err := controller.validateRule(&rule); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := controller.database.AddRule(&rule); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(201, rule)
}

func (controller *RulesController) validateRule(rule *rules.Rule) error {
	if rule.Name == "" {
		return &models.ValidationError{Message: "Name is required"}
	}

	_, err := rule.ReadAst()
	if err != nil {
		return &models.ValidationError{Message: err.Error()}
	}
	return nil
}
