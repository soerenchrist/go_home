package mqtt

import (
	"io"

	"github.com/gin-gonic/gin"
)

func (c *MqttController) publish(context *gin.Context) {
	topic := context.Param("topic")

	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.publishChannel <- Message{
		Topic:   topic,
		Payload: string(body),
	}

	context.JSON(200, gin.H{"success": true})
}
