package mqtt

import "github.com/gin-gonic/gin"

func NewMqttRouter(publishChannel PublishChannel) *gin.Engine {
	router := gin.Default()

	controller := MqttController{
		publishChannel: publishChannel,
	}

	v1 := router.Group("/api/v1")

	v1.POST("/publish/:topic", controller.publish)
	return router
}
