package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/soerenchrist/go_home/internal/config"
	"github.com/soerenchrist/go_home/internal/server"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: go_home -e <environment>")
		os.Exit(1)
	}

	flag.Parse()
	config.Init(*environment)

	server.Init()
}
