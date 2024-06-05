package main

import (
	"fmt"

	"github.com/its-a-feature/tiger/database"
	"github.com/its-a-feature/tiger/logging"
	"github.com/its-a-feature/tiger/rabbitmq"
	"github.com/its-a-feature/tiger/utils"
	"github.com/its-a-feature/tiger/webserver"
)

func main() {
	// initialize configuration based on .env and environment variables
	fmt.Print("Step 1/6 - Initializing utilities\n")
	utils.Initialize()
	// initialize logging
	fmt.Print("Step 2/6 - Initializing logging\n")
	logging.Initialize()
	// initialize the database connection
	// initialize database values if needed
	fmt.Print("Step 3/6 - Initializing database\n")
	database.Initialize()
	// initialize webserver
	fmt.Print("Step 4/6 - Initializing webserver\n")
	router := webserver.Initialize()
	// initialize all the rabbitmq connections
	// we will use rabbitmq to send messages to containers
	fmt.Print("Step 5/6 - Initializing rabbitmq\n")
	rabbitmq.Initialize()
	// start serving up API routes
	fmt.Print("Step 6/6 - Starting webserver\n")
	webserver.StartServer(router)
}
