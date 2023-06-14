package server

import (
	"fmt"

	"github.com/soerenchrist/mini_home/config"
	"github.com/soerenchrist/mini_home/db"
)

func Init() {
	database, err := db.NewDevicesDatabase("devices.db")
	if err != nil {
		panic(err)
	}

	config := config.GetConfig()
	r := NewRouter(database)

	port := config.GetString("server.port")
	host := config.GetString("server.host")

	addr := fmt.Sprintf("%s:%s", host, port)
	r.Run(addr)
}
