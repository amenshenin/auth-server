package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/amenshenin/auth-server.git/internal/config"
)

//"context"

func main() {
	//Глянуть что пихать в контекст и как его отменять
	//ctx = context.Background()

	//Sets flags
	configPath := flag.String("config-path", "/run/secrets/config", "Please set in the config-path flag the path to config file")
	flag.Parse()

	//Inits config
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err.Error())
	}
	fmt.Println(config)
}
