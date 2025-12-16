package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/amenshenin/auth-server.git/internal/config"
	"github.com/amenshenin/auth-server.git/internal/logger"
	"github.com/amenshenin/auth-server.git/internal/storage"
)

func main() {
	//Inits context
	ctx := context.Background()

	//Sets flags
	configPath := flag.String("config-path", "/run/secrets/config", "Please set in the config-path flag the path to config file")
	flag.Parse()

	//Inits config
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err.Error())
	}

	//Inits logger
	logger, err := logger.GetLogger(config)
	if err != nil {
		log.Fatalf("Error init logger: %s", err.Error())
	}
	logger.Info("Start service: init logger complete")

	//Inits DB
	db, err := storage.GetConnection(ctx, config, storage.PostgresProvider, logger)
	if err != nil {
		logger.Error("Error database connection", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Start service: getting database connection complete")
	defer db.Close()

	//Inits repo
	//Tnits service
	//Inits server & routs and starts it

	//Exits

	fmt.Println(config)
	fmt.Println(ctx)
	fmt.Println(db)
}
