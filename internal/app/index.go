package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) index(ctx *gin.Context) {
	devices, err := app.database.ListDevices()
	if err != nil {
		panic(err)
	}

	ctx.HTML(http.StatusOK, "index", gin.H{
		"title":   "Go Home!",
		"devices": devices,
	})
}
