package output

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WebsocketBinding struct {
	conn *websocket.Conn
}

var upgrader = websocket.Upgrader{}

func NewWebsocketBinding(router *gin.Engine) *WebsocketBinding {

	binding := &WebsocketBinding{}
	log.Info().Msg("Registering websocket handler")
	router.GET("/ws", binding.wsHandler)

	return binding
}

func (b *WebsocketBinding) Handle(value BindingValue) {
	log.Debug().Interface("value", value).
		Interface("conn", b.conn).Msg("Sending value to websocket")
	if b.conn == nil {
		return
	}

	b.conn.WriteJSON(value)
}

func (b *WebsocketBinding) wsHandler(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("Error upgrading websocket")
		return
	}

	b.conn = c

	defer c.Close()
	defer func() {
		b.conn = nil
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Err(err).Msg("Error reading message")
			break
		}
	}
}
