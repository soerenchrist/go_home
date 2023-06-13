package main

import (
	"github.com/gin-gonic/gin"
	"github.com/soerenchrist/mini_home/devices"
)

func main() {
	router := gin.Default()

	router.GET("/api/status", health)
	devices.MapEndpoints(router)
	router.Run()
}

func health(c *gin.Context) {
	c.Status(200)
}
