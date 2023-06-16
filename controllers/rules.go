package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/go_home/db"
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
	var rule rules.Rule
	if err := context.ShouldBindJSON(&rule); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := controller.database.AddRule(&rule); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(201, rule)
}
