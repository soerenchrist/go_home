package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) createRule(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "createRule", gin.H{})
}
