package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/soerenchrist/mini_home/config"
	"github.com/soerenchrist/mini_home/server"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: mini_home -e <environment>")
		os.Exit(1)
	}

	flag.Parse()
	config.Init(*environment)

	server.Init()
}
