package controllers

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/soerenchrist/mini_home/db"
	"github.com/soerenchrist/mini_home/models"
)

type CommandsController struct {
	database db.DevicesDatabase
}

func NewCommandsController(database db.DevicesDatabase) *CommandsController {
	return &CommandsController{database: database}
}

func (c *CommandsController) GetCommands(context *gin.Context) {
	deviceId := context.Param("deviceId")

	if _, err := c.database.GetDevice(deviceId); err != nil {
		context.JSON(404, gin.H{"error": "Device not found"})
		return
	}

	commands, err := c.database.ListCommands(deviceId)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, commands)
}

func (c *CommandsController) GetCommand(context *gin.Context) {
	deviceId := context.Param("deviceId")
	commandId := context.Param("commandId")

	if _, err := c.database.GetDevice(deviceId); err != nil {
		context.JSON(404, gin.H{"error": "Device not found"})
		return
	}

	command, err := c.database.GetCommand(deviceId, commandId)
	if err != nil {
		context.JSON(404, gin.H{"error": "Command not found"})
		return
	}

	context.JSON(200, command)
}

func (c *CommandsController) PostCommand(context *gin.Context) {
	deviceId := context.Param("deviceId")

	if _, err := c.database.GetDevice(deviceId); err != nil {
		context.JSON(404, gin.H{"error": "Device not found"})
		return
	}

	var request models.CreateCommandRequest
	if err := context.BindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := c.validateCommand(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	command := models.Command{
		ID:              uuid.NewString(),
		Name:            request.Name,
		DeviceID:        deviceId,
		PayloadTemplate: request.PayloadTemplate,
		Endpoint:        request.Endpoint,
		Method:          request.Method,
	}

	if err := c.database.AddCommand(command); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(201, command)
}

func (c *CommandsController) DeleteCommand(context *gin.Context) {
	deviceId := context.Param("deviceId")
	commandId := context.Param("commandId")

	if _, err := c.database.GetDevice(deviceId); err != nil {
		context.JSON(404, gin.H{"error": "Device not found"})
		return
	}

	if err := c.database.DeleteCommand(deviceId, commandId); err != nil {
		context.JSON(404, gin.H{"error": "Command not found"})
		return
	}

	context.Status(204)
}

func (c *CommandsController) InvokeCommand(context *gin.Context) {
	deviceId := context.Param("deviceId")
	commandId := context.Param("commandId")

	var params models.CommandParameters

	content_length := context.Request.Header["Content-Length"]
	if content_length != nil && content_length[0] != "0" {
		if err := context.BindJSON(&params); err != nil {
			log.Println("Failed to bind JSON", err)
		}
	}

	var device models.Device
	var err error
	if device, err = c.database.GetDevice(deviceId); err != nil {
		context.JSON(404, gin.H{"error": "Device not found"})
		return
	}

	command, err := c.database.GetCommand(deviceId, commandId)
	if err != nil {
		context.JSON(404, gin.H{"error": "Command not found"})
		return
	}

	var result models.InvocationResult
	res, err := command.Invoke(&device, &params)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	result.Response = string(body)
	result.StatusCode = res.StatusCode

	context.JSON(200, result)
}

func (c *CommandsController) validateCommand(command *models.CreateCommandRequest) error {
	if command.Name == "" {
		return &models.ValidationError{Message: "Name is required"}
	}

	if command.Endpoint == "" {
		return &models.ValidationError{Message: "Endpoint is required"}
	}
	if command.PayloadTemplate == "" {
		return &models.ValidationError{Message: "Payload template is required"}
	}

	methods := []string{"GET", "POST", "PUT", "DELETE"}
	if !contains(methods, command.Method) {
		return &models.ValidationError{Message: "Method must be one of GET, POST, PUT or DELETE"}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
