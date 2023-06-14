package controllers

import (
	"github.com/soerenchrist/mini_home/db"
	"github.com/soerenchrist/mini_home/models"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type DevicesController struct {
	database db.DevicesDatabase
}

func NewDevicesController(database db.DevicesDatabase) *DevicesController {
	return &DevicesController{database: database}
}

func (c *DevicesController) GetDevices(context *gin.Context) {
	devices, err := c.database.List()
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, devices)
}

func (c *DevicesController) PostDevice(context *gin.Context) {
	var request models.CreateDeviceRequest
	if err := context.BindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.validateDevice(request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	device := models.Device{
		ID:          uuid.NewString(),
		Name:        request.Name,
		LastReached: "Never",
	}

	err := c.database.Add(device)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, device)
}

func (c *DevicesController) validateDevice(device models.CreateDeviceRequest) error {
	if len(device.Name) < 3 {
		return &models.ValidationError{Message: "Name must be at least 3 characters long"}
	}

	return nil
}
