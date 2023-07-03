package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) createDevice(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "createDevice", gin.H{})
}
