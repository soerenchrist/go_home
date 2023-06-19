package controllers

import (
	"fmt"

	"github.com/soerenchrist/go_home/models"

	"github.com/gin-gonic/gin"
)

type DevicesDatabase interface {
	ListDevices() ([]models.Device, error)
	GetDevice(deviceId string) (*models.Device, error)
	AddDevice(device *models.Device) error
	DeleteDevice(deviceId string) error
}

type DevicesController struct {
	database DevicesDatabase
}

func NewDevicesController(database DevicesDatabase) *DevicesController {
	return &DevicesController{database: database}
}

func (c *DevicesController) GetDevices(context *gin.Context) {
	devices, err := c.database.ListDevices()
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

	_, err := c.database.GetDevice(request.Id)
	if err == nil {
		context.JSON(400, gin.H{"error": fmt.Sprintf("Device with id %s does already exist", request.Id)})
		return
	}

	if err := c.validateDevice(request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	device := models.Device{
		ID:   request.Id,
		Name: request.Name,
	}

	err = c.database.AddDevice(&device)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(201, device)
}

func (c *DevicesController) GetDevice(context *gin.Context) {
	id := context.Param("deviceId")

	device, err := c.database.GetDevice(id)
	if err != nil {
		context.JSON(404, gin.H{"error": "Device not found"})
		return
	}
	context.JSON(200, device)
}

func (c *DevicesController) DeleteDevice(context *gin.Context) {
	id := context.Param("deviceId")

	err := c.database.DeleteDevice(id)

	if notFound, isOk := err.(*models.NotFoundError); isOk {
		context.JSON(404, gin.H{"error": notFound.Error()})
		return
	}

	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.Status(204)
}

func (c *DevicesController) validateDevice(device models.CreateDeviceRequest) error {
	if len(device.Name) < 3 {
		return &models.ValidationError{Message: "Name must be at least 3 characters long"}
	}

	return nil
}
