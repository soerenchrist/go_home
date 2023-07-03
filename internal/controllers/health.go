package controllers

import "github.com/gin-gonic/gin"

type HealthController struct {
}

func (c *HealthController) Status(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": "ok",
	})
}
