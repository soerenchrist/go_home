package devices

import (
	"github.com/soerenchrist/mini_home/errors"
	"github.com/soerenchrist/mini_home/persistence"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

var database persistence.Database[Device]

func getDevices(context *gin.Context) {
	devices, err := database.List()
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, devices)
}

func postDevice(context *gin.Context) {
	var request CreateDeviceRequest
	if err := context.BindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validateDevice(request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	device := Device{
		ID:          uuid.NewString(),
		Name:        request.Name,
		LastReached: "Never",
	}

	err := database.Add(device)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, device)
}

func validateDevice(device CreateDeviceRequest) error {
	if len(device.Name) < 3 {
		return &errors.ValidationError{Message: "Name must be at least 3 characters long"}
	}

	return nil
}

func MapEndpoints(router *gin.Engine) {
	var err error
	database, err = NewSqliteDatabase("devices.db")
	if err != nil {
		panic(err)
	}
	router.GET("/api/devices", getDevices)
	router.POST("/api/devices", postDevice)
}
