package server

import (
	"html/template"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	frontend "github.com/soerenchrist/go_home/internal/app"
	"github.com/soerenchrist/go_home/internal/command"
	"github.com/soerenchrist/go_home/internal/db"
	"github.com/soerenchrist/go_home/internal/device"
	"github.com/soerenchrist/go_home/internal/rules"
	"github.com/soerenchrist/go_home/internal/sensor"
	"github.com/soerenchrist/go_home/internal/value"
	"github.com/soerenchrist/go_home/pkg/output"
)

func NewRouter(database db.Database, outputBindings *output.OutputBindingsManager) *gin.Engine {
	router := gin.New()
	router.Use(DefaultStructuredLogger())
	router.Use(gin.Recovery())

	app := frontend.NewApp(router, database)
	app.ServeHtml()

	devicesController := device.NewController(database)
	sensorsController := sensor.NewController(database)
	sensorValuesController := value.NewController(database, outputBindings)
	commandsController := command.NewController(database)
	rulesController := rules.NewController(database)

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", health)

	v1.GET("/devices", devicesController.GetDevices)
	v1.GET("/devices/:deviceId", devicesController.GetDevice)
	v1.POST("/devices", devicesController.PostDevice)
	v1.DELETE("/devices/:deviceId", devicesController.DeleteDevice)

	v1.GET("/devices/:deviceId/sensors", sensorsController.GetSensors)
	v1.POST("/devices/:deviceId/sensors", sensorsController.PostSensor)
	v1.GET("/devices/:deviceId/sensors/:sensorId", sensorsController.GetSensor)
	v1.DELETE("/devices/:deviceId/sensors/:sensorId", sensorsController.DeleteSensor)

	v1.POST("/devices/:deviceId/sensors/:sensorId/values", sensorValuesController.PostSensorValue)
	v1.GET("/devices/:deviceId/sensors/:sensorId/values", sensorValuesController.GetSensorValues)
	v1.GET("/devices/:deviceId/sensors/:sensorId/current", sensorValuesController.GetCurrentSensorValue)

	v1.GET("/devices/:deviceId/commands", commandsController.GetCommands)
	v1.GET("/devices/:deviceId/commands/:commandId", commandsController.GetCommand)
	v1.POST("/devices/:deviceId/commands", commandsController.PostCommand)
	v1.POST("/devices/:deviceId/commands/:commandId/invoke", commandsController.InvokeCommand)
	v1.DELETE("/devices/:deviceId/commands/:commandId", commandsController.DeleteCommand)

	v1.GET("/rules", rulesController.ListRules)
	v1.POST("/rules", rulesController.PostRule)

	router.POST("/echo", echo)
	router.GET("/websocket", websocketPage)
	return router
}

func websocketPage(ctx *gin.Context) {
	websocketTemplate.Execute(ctx.Writer, "ws://"+ctx.Request.Host+"/ws")
}

func echo(context *gin.Context) {
	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		log.Error().Msgf("Error reading body: %s", err.Error())
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.Status(200)
	context.Header("Content-Type", "application/json")
	context.Writer.Write(body)
}

func health(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": "ok",
	})
}

var websocketTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
