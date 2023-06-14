package server

import (
	"fmt"

	"github.com/soerenchrist/mini_home/config"
	"github.com/soerenchrist/mini_home/db"
)

func Init() {

	config := config.GetConfig()
	databasePath := config.GetString("database.path")
	database, err := db.NewDevicesDatabase(databasePath)
	if err != nil {
		panic(err)
	}
	r := NewRouter(database)

	port := config.GetString("server.port")
	host := config.GetString("server.host")

	addr := fmt.Sprintf("%s:%s", host, port)
	r.Run(addr)
}
