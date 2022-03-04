package api

import (
	"fmt"
	"os"
	"time"

	"github.com/Hsmnasiri/http_monitoring/server/api/controllers"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error getting env!")
		time.Sleep(5 * time.Second)
		server.Initialize("postgres", "postgres", "postgres", "5432", "pgdb", "postgres")
	} else {
		fmt.Println("We are getting the env values!")
		server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	}

	//seed.Load(server.DB)

	server.Run(":8080")

}
