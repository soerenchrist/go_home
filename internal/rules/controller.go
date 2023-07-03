package rules

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/internal/errors"
)

type RulesController struct {
	database RulesDatabase
}

func NewController(database RulesDatabase) *RulesController {
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
	var request CreateRuleRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rule := Rule{
		Name: request.Name,
		When: WhenExpression(request.When),
		Then: ThenExpression(request.Then),
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

func (controller *RulesController) validateRule(rule *Rule) error {
	if rule.Name == "" {
		return &errors.ValidationError{Message: "Name is required"}
	}

	_, err := rule.ReadConditionAst()
	if err != nil {
		return &errors.ValidationError{Message: err.Error()}
	}

	_, err = rule.ReadAction()
	if err != nil {
		log.Printf("Error: %v", err)
		return &errors.ValidationError{Message: err.Error()}
	}
	return nil
}
